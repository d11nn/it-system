package context_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/Alonza0314/it-system/controller/backend/internal/context"
	"github.com/Alonza0314/it-system/controller/backend/model"
)

func TestCreateTaskStoresLibraryPrList(t *testing.T) {
	tempDir := t.TempDir()
	localCtx := context.NewItContext(filepath.Join(tempDir, "test.db"), filepath.Join(tempDir, "log"), 20, 30*time.Second, false, "", nil)
	defer func() {
		if err := context.ReleaseItContext(localCtx); err != nil {
			t.Fatalf("release context: %v", err)
		}
	}()

	err := localCtx.CreateTask(
		"tester",
		123,
		[]string{"TestRegistration"},
		[]model.NfPr{{NfName: "amf", PR: 204}},
		[]model.LibraryPr{{RepoName: "openapi", PR: 67}},
	)
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	tasks := localCtx.GetPendingTasks()
	if len(tasks) != 1 {
		t.Fatalf("pending task length = %d, want 1", len(tasks))
	}

	task, err := localCtx.GetTask(tasks[0].Id)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}

	libraryPrList := task.LibraryPrList()
	if len(libraryPrList) != 1 {
		t.Fatalf("libraryPrList length = %d, want 1", len(libraryPrList))
	}
	if libraryPrList[0].RepoName() != "openapi" || libraryPrList[0].PR() != 67 {
		t.Fatalf("libraryPrList[0] = %s #%d, want openapi #67", libraryPrList[0].RepoName(), libraryPrList[0].PR())
	}
}
