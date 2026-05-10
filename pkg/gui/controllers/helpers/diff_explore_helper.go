package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type DiffExploreHelper struct {
	c *HelperCommon
}

func NewDiffExploreHelper(c *HelperCommon) *DiffExploreHelper {
	return &DiffExploreHelper{c: c}
}

func (self *DiffExploreHelper) RefreshDiffExplorePanel(focusOpts types.OnFocusOpts) {
	if self.c.Context().CurrentStatic().GetKey() != self.c.Contexts().DiffExplore.GetKey() {
		return
	}

	mainContext := self.c.Contexts().DiffExplore

	var file *models.File
	node := self.c.Contexts().Files.GetSelected()
	if node != nil {
		file = node.File
	}

	if file == nil || (!file.HasUnstagedChanges && !file.HasStagedChanges) {
		self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
		return
	}

	diff := self.diffForFile(file)

	mainContext.GetMutex().Lock()
	defer mainContext.GetMutex().Unlock()

	selectedLineIdx := -1
	if focusOpts.ClickedViewLineIdx > 0 {
		selectedLineIdx = focusOpts.ClickedViewLineIdx
	}

	hunkMode := self.c.UserConfig().Gui.UseHunkModeInStagingView
	mainContext.SetState(
		patch_exploring.NewState(diff, selectedLineIdx, mainContext.GetView(), mainContext.GetState(), hunkMode),
	)

	state := mainContext.GetState()
	if state == nil {
		self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
		return
	}

	content := mainContext.GetContentToRender()
	mainContext.FocusSelection()

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().DiffExplore,
		Main: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(content),
			Title: self.c.Tr.DiffTitle,
		},
	})
}

func (self *DiffExploreHelper) diffForFile(file *models.File) string {
	worktree := self.c.Git().WorkingTree
	if !file.GetIsTracked() && !file.HasStagedChanges {
		return worktree.WorktreeFileDiff(file, true, false)
	}
	return worktree.WorktreeFileCombinedDiff(file, true)
}
