package bot

import (
	"fmt"
	"strings"
)

func (b *bot) handleStartSession(chatID int64, sessionName string) string {
	var response string
	session, err := b.storage.StartSession(sessionName, chatID)
	if err != nil {
		response = fmt.Sprintf("Failed to create Session, invalid session name")
		log.Errorf("Failed to create Session, invalid session arguments: %s", sessionName)
		return response
	}
	b.ActiveSessions[chatID] = session
	log.Infof("Starting Session %v\n", session)
	response = fmt.Sprintf("Starting \\#session: \n \U0001F3F7  *_\\#%s_*", sessionName)
	return response
}

func (b *bot) handleRenameSession(chatID int64, name string) string {
	activeSession, ok := b.ActiveSessions[chatID]
	if ok == false {
		response := fmt.Sprintf("Failed to rename Session, no active session found")
		log.Errorf("Failed to rename Session session not found: %d", chatID)
		return response
	}

	oldName := activeSession.Name
	err := b.storage.RenameSession(activeSession, name)
	if err != nil {
		response := fmt.Sprintf("Failed to rename Session _\\#%s_, invalid session name", oldName)
		log.Errorf("Failed to rename Session _\\#%s_, invalid session arguments: %s", oldName, name)
		return response
	}
	return fmt.Sprintf(label+" \\#session renamed: \n _\\#%s_ "+rigthArrow+" *_\\#%s_*", oldName, name)
}

func (b *bot) handleEndSession(chatID int64) string {
	activeSession, ok := b.ActiveSessions[chatID]
	var response strings.Builder
	if ok == false {
		response.WriteString(fmt.Sprintf("Failed to rename Session, no active session found"))
		log.Errorf("Failed to rename Session session not found: %d", chatID)
		return response.String()
	}
	log.Info("Closing Session: ", activeSession)
	response.WriteString(fmt.Sprintf(label+" \\#session _\\#%s_ Finished\n", activeSession.Name))
	// before closing the session get the summay
	summary, err := b.storage.Summary(activeSession)
	if err != nil {
		log.Error(err)
	}
	if len(summary) > 0 {
		response.WriteString("Player             \\| Expression \\| Rolls \\| Min  \\| Max  \\| Average       \n")
		for _, sl := range summary {
			response.WriteString(sl.String())
		}
		response.WriteString("\n")
	}
	b.storage.EndSession(activeSession)
	delete(b.ActiveSessions, chatID)
	log.Info(response.String())
	return response.String()
}
