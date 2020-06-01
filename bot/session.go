package bot

import "fmt"

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
	response = fmt.Sprintf("Starting Session: \n \U0001F3F7  *_\\#%s_*", sessionName)
	return response
}

func (b *bot) handleRenameSession(chatID int64, name string) string {
	activeSession := b.ActiveSessions[chatID]
	if activeSession == nil {
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
	return fmt.Sprintf(label+" Session renamed: \n _\\#%s_ "+rigthArrow+" *_\\#%s_*", oldName, name)
}

func (b *bot) handleEndSession(chatID int64) string {
	activeSession := b.ActiveSessions[chatID]
	if activeSession == nil {
		response := fmt.Sprintf("Failed to rename Session, no active session found")
		log.Errorf("Failed to rename Session session not found: %d", chatID)
		return response
	}
	log.Info("Closing Session: ", activeSession)
	b.storage.EndSession(activeSession)
	b.ActiveSessions[chatID] = nil
	response := fmt.Sprintf(label+" Session _\\#%s_ Finished", activeSession.Name)
	log.Info(response)
	return response
}
