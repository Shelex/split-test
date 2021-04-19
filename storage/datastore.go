package storage

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/Shelex/split-specs/entities"
)

type DataStore struct {
	Client *datastore.Client
}

func NewDataStore() (Storage, error) {
	store := DataStore{}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "test-splitter")

	if err != nil {
		return store, err
	}

	//defer client.Close()

	store.Client = client

	DB = store

	return store, nil
}

func (d DataStore) CreateProject(project entities.Project) error {
	return errors.New("not implemented")
}
func (d DataStore) CreateSession(projectName string, sessionID string, specs []entities.Spec) (*entities.Session, error) {
	return nil, errors.New("not implemented")
}
func (d DataStore) AttachSessionToProject(projectName string, sessionID string) error {
	return errors.New("not implemented")
}
func (d DataStore) GetProjectLatestSession(projectName string) (*entities.Session, error) {
	return nil, errors.New("not implemented")
}
func (d DataStore) SetProjectLatestSession(projectName string, sessionID string) error {
	return errors.New("not implemented")
}
func (d DataStore) GetFullProjectByName(name string) (entities.ProjectFull, error) {
	return entities.ProjectFull{}, errors.New("not implemented")
}
func (d DataStore) StartSpec(sessionID string, machineID string, specName string) error {
	return errors.New("not implemented")
}
func (d DataStore) EndSpec(sessionID string, machineID string) error {
	return errors.New("not implemented")
}
func (d DataStore) GetSession(sessionID string) (entities.Session, error) {
	return entities.Session{}, errors.New("not implemented")
}
func (d DataStore) EndSession(sessionID string) error {
	return errors.New("not implemented")
}

func (d DataStore) CreateUser(entities.User) error {
	return errors.New("not implemented")
}

func (d DataStore) GetUserByUsername(username string) (*entities.User, error) {
	return nil, errors.New("not implemented")
}

func (d DataStore) GetUserProjectIDByName(userID string, projectName string) (string, error) {
	return "", errors.New("not implemented")
}

func (d DataStore) GetProjectByID(ID string) (*entities.Project, error) {
	return nil, errors.New("not implemented")
}

func (d DataStore) AttachProjectToUser(userID string, projectID string) error {
	return errors.New("not implemented")
}

func (d DataStore) GetUserProjects(userID string) ([]string, error) {
	return []string{}, errors.New("not implemented")
}

func (d DataStore) UpdatePassword(userID string, newPassword string) error {
	return errors.New("not implemented")
}
