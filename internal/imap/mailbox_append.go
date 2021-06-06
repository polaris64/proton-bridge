// Copyright (c) 2021 Proton Technologies AG
//
// This file is part of ProtonMail Bridge.
//
// ProtonMail Bridge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// ProtonMail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ProtonMail Bridge.  If not, see <https://www.gnu.org/licenses/>.

package imap

import (
	"context"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/proton-bridge/internal/imap/uidplus"
	"github.com/ProtonMail/proton-bridge/pkg/message"
	"github.com/ProtonMail/proton-bridge/pkg/pmapi"
	"github.com/emersion/go-imap"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CreateMessage appends a new message to this mailbox. The \Recent flag will
// be added regardless of whether flags is empty or not. If date is nil, the
// current time will be used.
//
// If the Backend implements Updater, it must notify the client immediately
// via a mailbox update.
func (im *imapMailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return im.logCommand(func() error {
		return im.createMessage(flags, date, body)
	}, "APPEND", flags, date)
}

func (im *imapMailbox) createMessage(flags []string, date time.Time, body imap.Literal) error { //nolint[funlen]
	// Called from go-imap in goroutines - we need to handle panics for each function.
	defer im.panicHandler.HandlePanic()

	m, _, _, readers, err := message.Parse(body)
	if err != nil {
		return err
	}

	apiCtx := context.Background()
	labelIDs := getLabelIDsFromKeywords(m, im.user.client(), apiCtx)

	logrus.Info("Append the following label IDs to m.LabelIDs: -")
	logrus.Info(labelIDs)
	m.LabelIDs = append(m.LabelIDs, labelIDs...)

	addr := im.storeAddress.APIAddress()
	if addr == nil {
		return errors.New("no available address for encryption")
	}
	m.AddressID = addr.ID

	kr, err := im.user.client().KeyRingForAddressID(addr.ID)
	if err != nil {
		return err
	}

	// Handle imported messages which have no "Sender" address.
	// This sometimes occurs with outlook which reports errors as imported emails or for drafts.
	if m.Sender == nil {
		im.log.Warning("Append: Missing email sender. Will use main address")
		m.Sender = &mail.Address{
			Name:    "",
			Address: addr.Email,
		}
	}

	// "Drafts" needs to call special API routes.
	// Clients always append the whole message again and remove the old one.
	if im.storeMailbox.LabelID() == pmapi.DraftLabel {
		// Sender address needs to be sanitised (drafts need to match cases exactly).
		m.Sender.Address = pmapi.ConstructAddress(m.Sender.Address, addr.Email)

		draft, _, err := im.user.storeUser.CreateDraft(kr, m, readers, "", "", "")
		if err != nil {
			return errors.Wrap(err, "failed to create draft")
		}

		targetSeq := im.storeMailbox.GetUIDList([]string{draft.ID})
		return uidplus.AppendResponse(im.storeMailbox.UIDValidity(), targetSeq)
	}

	// We need to make sure this is an import, and not a sent message from this account
	// (sent messages from the account will be added by the event loop).
	if im.storeMailbox.LabelID() == pmapi.SentLabel {
		sanitizedSender := pmapi.SanitizeEmail(m.Sender.Address)

		// Check whether this message was sent by a bridge user.
		user, err := im.user.backend.bridge.GetUser(sanitizedSender)
		if err == nil && user.ID() == im.storeUser.UserID() {
			logEntry := im.log.WithField("addr", sanitizedSender).WithField("extID", m.Header.Get("Message-Id"))

			// If we find the message in the store already, we can skip importing it.
			if foundUID := im.storeMailbox.GetUIDByHeader(&m.Header); foundUID != uint32(0) {
				logEntry.Info("Ignoring APPEND of duplicate to Sent folder")
				return uidplus.AppendResponse(im.storeMailbox.UIDValidity(), &uidplus.OrderedSeq{foundUID})
			}

			// We didn't find the message in the store, so we are currently sending it.
			logEntry.WithField("time", date).Info("No matching UID, continuing APPEND to Sent")
		}
	}

	message.ParseFlags(m, flags)
	if !date.IsZero() {
		m.Time = date.Unix()
	}

	internalID := m.Header.Get("X-Pm-Internal-Id")
	references := m.Header.Get("References")
	referenceList := strings.Fields(references)

	// In case there is a mail client which corrupts headers, try
	// "References" too.
	if internalID == "" && len(referenceList) > 0 {
		lastReference := referenceList[len(referenceList)-1]
		match := pmapi.RxInternalReferenceFormat.FindStringSubmatch(lastReference)
		if len(match) == 2 {
			internalID = match[1]
		}
	}

	im.user.appendExpungeLock.Lock()
	defer im.user.appendExpungeLock.Unlock()

	// Avoid appending a message which is already on the server. Apply the
	// new label instead. This always happens with Outlook (it uses APPEND
	// instead of COPY).
	if internalID != "" {
		// Check to see if this belongs to a different address in split mode or another ProtonMail account.
		msg, err := im.storeMailbox.GetMessage(internalID)
		if err == nil && (im.user.user.IsCombinedAddressMode() || (im.storeAddress.AddressID() == msg.Message().AddressID)) {
			IDs := []string{internalID}

			// See the comment bellow.
			if msg.IsMarkedDeleted() {
				if err := im.storeMailbox.MarkMessagesUndeleted(IDs); err != nil {
					log.WithError(err).Error("Failed to undelete re-imported internal message")
				}
			}

			err = im.storeMailbox.LabelMessages(IDs)
			if err != nil {
				return err
			}

			// Add labels to messages that already exist
			logrus.Info(fmt.Sprintf("Message already exists (ID: %s), adding labels...", internalID))
			addCustomLabels(labelIDs, internalID, im.user.client(), apiCtx)

			targetSeq := im.storeMailbox.GetUIDList(IDs)
			return uidplus.AppendResponse(im.storeMailbox.UIDValidity(), targetSeq)
		}
	}

	im.log.Info("Importing external message")
	if err := im.importMessage(m, readers, kr, labelIDs); err != nil {
		im.log.Error("Import failed: ", err)
		return err
	}

	// IMAP clients can move message to local folder (setting \Deleted flag)
	// and then move it back (IMAP client does not remember the message,
	// so instead removing the flag it imports duplicate message).
	// Regular IMAP server would keep the message twice and later EXPUNGE would
	// not delete the message (EXPUNGE would delete the original message and
	// the new duplicate one would stay). API detects duplicates; therefore
	// we need to remove \Deleted flag if IMAP client re-imports.
	msg, err := im.storeMailbox.GetMessage(m.ID)
	if err == nil && msg.IsMarkedDeleted() {
		if err := im.storeMailbox.MarkMessagesUndeleted([]string{m.ID}); err != nil {
			log.WithError(err).Error("Failed to undelete re-imported message")
		}
	}

	targetSeq := im.storeMailbox.GetUIDList([]string{m.ID})
	return uidplus.AppendResponse(im.storeMailbox.UIDValidity(), targetSeq)
}

func (im *imapMailbox) importMessage(m *pmapi.Message, readers []io.Reader, kr *crypto.KeyRing, labelIDs []string) (err error) {
	body, err := message.BuildEncrypted(m, readers, kr)
	if err != nil {
		return err
	}

	labels := []string{}
	for _, l := range m.LabelIDs {
		if l == pmapi.StarredLabel {
			labels = append(labels, pmapi.StarredLabel)
		}
	}

	// Append extra labels
	for _, l := range labelIDs {
		labels = append(labels, l)
	}

	return im.storeMailbox.ImportMessage(m, body, labels)
}

func getLabelIDsFromKeywords(m *pmapi.Message, api pmapi.Client, ctx context.Context) []string {
	var keywordsHeader string
	keywordsHeader = m.Header.Get("X-Keywords")
	if len(keywordsHeader) == 0 {
	    keywordsHeader = m.Header.Get("X-Label")
	}
	if len(keywordsHeader) == 0 {
	    keywordsHeader = m.Header.Get("Keywords")
	}
	if len(keywordsHeader) == 0 {
		return []string{}
	}
	logrus.Info(fmt.Sprintf("Found keywords: %s", keywordsHeader))

	// Convert all keywords (label names) to IDs (creating if necessary)
	labelNames := strings.Split(keywordsHeader, ",")
	var labels []*pmapi.Label
	labels = api.GetLabelCache()
	if len(labels) == 0 {
		labels2, err := api.ListLabels(ctx)
		if err != nil {
			logrus.Error(err)
		} else {
			labels = labels2
		}
	}
	labelIDs := []string{}
	found := false
	for _, keyword := range labelNames {
		keyword = strings.TrimSpace(keyword)
		found = false
		for _, label := range labels {
			if label.Name == keyword {
				labelIDs = append(labelIDs, label.ID)
				found = true
				break
			}
		}

		if !found {
			// Create new Label and append ID
			logrus.Warn(fmt.Sprintf("Label \"%s\" not found, creating...", keyword))
			label := pmapi.Label{
				Name:      keyword,
				Path:      keyword,
				Color:     pmapi.LeastUsedColor(pmapi.LabelColors),
				Display:   0,
				Exclusive: false,
				Type:      1,
				Notify:    false,
			}
			newLabel, err := api.CreateLabel(ctx, &label)
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.Info(fmt.Sprintf("Created new label \"%s\" (ID: %s)", keyword, newLabel.ID))
				labelIDs = append(labelIDs, newLabel.ID)
			}
		}
	}
	return labelIDs
}


func addCustomLabels(labelIDs []string, internalID string, api pmapi.Client, ctx context.Context) {
	if len(labelIDs) == 0 {
		return
	}

	apiMsg, err := api.GetMessage(ctx, internalID)

	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("The message currently has the following labels")
		logrus.Info(apiMsg.LabelIDs)

		// Remove existing Message Labels which are not in labelNames
		removeLabelIDs := []string{}
		for _, exLabelID := range apiMsg.LabelIDs {
			if !pmapi.IsSystemLabel(exLabelID) {
				found := false
				for _, newLabelID := range labelIDs {
					if newLabelID == exLabelID {
						found = true
						break
					}
				}
				if !found {
					removeLabelIDs = append(removeLabelIDs, exLabelID)
				}
			}
		}

		for _, remLabelID := range removeLabelIDs {
			logrus.Info(fmt.Sprintf("Removing label ID %s...", remLabelID))
			api.UnlabelMessages(ctx, []string{internalID}, remLabelID)
		}
	}

	for _, labelID := range labelIDs {
		logrus.Info(fmt.Sprintf("Adding label ID %s...", labelID))
		api.LabelMessages(ctx, []string{internalID}, labelID)
	}
}
