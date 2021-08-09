package storage

import (
	"fmt"
	"time"

	"github.com/Shelex/split-specs/entities"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type InMem struct {
	sessions map[string]*entities.Session
	projects map[string]*entities.Project
	users    map[string]*entities.User
	specs    map[string]*entities.Spec
}

func NewInMemStorage() (Storage, error) {
	DB = &InMem{
		sessions: map[string]*entities.Session{},
		projects: map[string]*entities.Project{},
		users:    map[string]*entities.User{},
		specs:    map[string]*entities.Spec{},
	}
	return DB, nil
}

func (i *InMem) CreateUser(userInput entities.User) error {
	i.users[userInput.ID] = &userInput
	return nil
}

func (i *InMem) UpdatePassword(userID string, newPassword string) error {
	i.users[userID].Password = newPassword
	return nil
}

func (i *InMem) GetUserByEmail(email string) (*entities.User, error) {
	for _, user := range i.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (i *InMem) GetUserProjects(userID string) ([]string, error) {
	return i.users[userID].ProjectIDs, nil
}

func (i *InMem) GetUserProjectIDByName(userID string, projectName string) (string, error) {
	user, ok := i.users[userID]
	if !ok {
		return "", fmt.Errorf("user not found")
	}
	for _, ID := range user.ProjectIDs {
		project, err := i.GetProjectByID(ID)
		if err == nil && project.Name == projectName {
			return project.ID, nil
		}
	}
	return "", ErrProjectNotFound
}

func (i *InMem) GetProjectByID(ID string) (*entities.Project, error) {
	project, ok := i.projects[ID]
	if !ok {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

func (i *InMem) CreateProject(project entities.Project) error {
	i.projects[project.ID] = &project
	return nil
}

func (i *InMem) AttachProjectToUser(userID string, projectID string) error {
	i.users[userID].ProjectIDs = append(i.users[userID].ProjectIDs, projectID)
	return nil
}

func (i *InMem) CreateSession(projectID string, sessionID string, specs []entities.Spec) (*entities.Session, error) {
	if _, ok := i.sessions[sessionID]; ok {
		return nil, fmt.Errorf("[repository]: session id already in use for project %s", projectID)
	}

	err := i.CreateSpecs(sessionID, specs)
	if err != nil {
		return nil, fmt.Errorf("failed to create specs")
	}

	session := &entities.Session{
		ID:        sessionID,
		ProjectID: projectID,
	}

	i.sessions[sessionID] = session
	return session, nil
}

func (i *InMem) GetProjectLatestSessions(projectID string, limit int) ([]*entities.Session, error) {
	var sessions []*entities.Session

	projectSessions, err := i.GetProjectSessions(projectID)
	if err != nil {
		return nil, err
	}

	for _, projectSession := range projectSessions {
		session, ok := i.sessions[projectSession.ID]
		if !ok {
			return nil, ErrSessionNotFound
		}
		if session.End != 0 {
			sessions = append(sessions, session)
		}

		if len(sessions) >= limit {
			return sessions, nil
		}
	}
	return sessions, nil
}

func (i *InMem) GetProjectLatestSession(projectID string) (*entities.Session, error) {
	project, ok := i.projects[projectID]
	if !ok {
		return nil, ErrProjectNotFound
	}

	latestSession, ok := i.sessions[project.LatestSession]
	if !ok {
		return nil, fmt.Errorf("latest session for project %s not found", projectID)
	}

	return latestSession, nil
}

func (i *InMem) SetProjectLatestSession(projectID string, sessionID string) error {
	_, ok := i.projects[projectID]
	if !ok {
		return ErrProjectNotFound
	}
	i.projects[projectID].LatestSession = sessionID
	return nil
}

func (i *InMem) GetSession(sessionID string) (entities.Session, error) {
	var empty entities.Session
	session, ok := i.sessions[sessionID]
	if !ok {
		return empty, fmt.Errorf("session %s not found", sessionID)

	}
	return *session, nil
}

func (i *InMem) StartSpec(sessionID string, machineID string, specName string) error {
	session, err := i.GetSession(sessionID)
	if err != nil {
		return err
	}

	specs, err := i.GetSpecs(sessionID)
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if spec.FilePath == specName {
			if session.Start == 0 {
				i.sessions[sessionID].Start = time.Now().Unix()
			}
			i.specs[spec.ID].Start = time.Now().Unix()
			i.specs[spec.ID].AssignedTo = machineID
			return nil
		}
	}
	return nil
}

func (i *InMem) EndSpec(sessionID string, machineID string, isPassed bool) error {
	specs, err := i.GetSpecs(sessionID)
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if spec.End == 0 && spec.Start != 0 && spec.AssignedTo == machineID {
			i.specs[spec.ID].End = time.Now().Unix()
			i.specs[spec.ID].EstimatedDuration = i.specs[spec.ID].End - i.specs[spec.ID].Start
			i.specs[spec.ID].Passed = isPassed
			return nil
		}
	}
	return nil
}

func (i *InMem) EndSession(sessionID string) error {
	session, err := i.GetSession(sessionID)
	if err != nil {
		return err
	}

	i.sessions[sessionID].End = time.Now().Unix()

	if err := i.SetProjectLatestSession(session.ProjectID, sessionID); err != nil {
		return err
	}

	return nil
}

func (i *InMem) CreateSpecs(sessionID string, specs []entities.Spec) error {
	for _, spec := range specs {
		id, _ := gonanoid.New()
		spec.ID = id
		spec.SessionID = sessionID
		i.specs[spec.ID] = &spec
	}
	return nil
}

func (i *InMem) GetSpecs(sessionID string) ([]entities.Spec, error) {
	var specs []entities.Spec

	for _, spec := range i.specs {
		if spec.SessionID == sessionID {
			specs = append(specs, *spec)
		}
	}

	return specs, nil
}

func (i *InMem) DeleteProject(email string, projectID string) error {
	projectSessions, err := i.GetProjectSessions(projectID)
	if err != nil {
		return err
	}

	for _, session := range projectSessions {
		err := i.DeleteSession(email, session.ID)
		if err != nil {
			return err
		}
	}
	delete(i.projects, projectID)

	for _, user := range i.users {
		userHasProject, index := contains(user.ProjectIDs, projectID)
		if userHasProject {
			i.users[user.ID].ProjectIDs = remove(i.users[user.ID].ProjectIDs, index)
		}
	}
	return nil
}

func (i *InMem) DeleteSession(email string, sessionID string) error {
	session, ok := i.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	for _, spec := range i.specs {
		if spec.SessionID == sessionID {
			delete(i.specs, spec.ID)
		}
	}

	delete(i.sessions, session.ID)
	return nil
}

func (i *InMem) GetProjectSessions(projectID string) ([]entities.SessionWithSpecs, error) {
	var sessions []entities.SessionWithSpecs
	for _, session := range i.sessions {
		if session.ProjectID == projectID {
			sessionWithSpecs, err := i.GetSessionWithSpecs(session.ID)
			if err != nil {
				return sessions, err
			}
			sessions = append(sessions, sessionWithSpecs)
		}
	}

	return sessions, nil
}

func (i *InMem) GetSessionWithSpecs(sessionID string) (entities.SessionWithSpecs, error) {
	var empty entities.SessionWithSpecs
	session, err := i.GetSession(sessionID)
	if err != nil {
		return empty, err
	}

	specs, err := i.GetSpecs(sessionID)
	if err != nil {
		return empty, err
	}

	return entities.SessionWithSpecs{
		ID:        session.ID,
		ProjectID: session.ProjectID,
		Start:     session.Start,
		End:       session.End,
		Specs:     specs,
	}, nil
}

func contains(input []string, query string) (bool, int) {
	for index, item := range input {
		if item == query {
			return true, index
		}
	}
	return false, -1
}

func remove(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
