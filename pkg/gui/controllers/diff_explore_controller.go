package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type DiffExploreController struct {
	baseController
	c                 *ControllerCommon
	context           types.IPatchExplorerContext
	prevExtrasVisible bool
}

func NewDiffExploreController(c *ControllerCommon, context types.IPatchExplorerContext) *DiffExploreController {
	return &DiffExploreController{
		baseController: baseController{},
		c:              c,
		context:        context,
	}
}

func (self *DiffExploreController) Context() types.Context {
	return self.context
}

func (self *DiffExploreController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Edit),
			Handler:         self.editAtCursor,
			Description:     self.c.Tr.EditFile,
			Tooltip:         self.c.Tr.EditFileTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.OpenFile),
			Handler:     self.openFile,
			Description: self.c.Tr.OpenFile,
			Tooltip:     self.c.Tr.OpenFileTooltip,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Return),
			Handler:         self.escape,
			Description:     self.c.Tr.ReturnToFilesPanel,
			DisplayOnScreen: true,
		},
	}
}

func (self *DiffExploreController) selectedFilePath() string {
	node := self.c.Contexts().Files.GetSelected()
	if node == nil || node.File == nil {
		return ""
	}
	return node.File.GetPath()
}

func (self *DiffExploreController) editAtCursor() error {
	self.context.GetMutex().Lock()
	defer self.context.GetMutex().Unlock()

	path := self.selectedFilePath()
	if path == "" {
		return nil
	}

	state := self.context.GetState()
	if state == nil {
		return nil
	}

	lineNumber := state.CurrentLineNumber()
	lineNumber = self.c.Helpers().Diff.AdjustLineNumber(path, lineNumber, self.context.GetViewName())
	return self.c.Helpers().Files.EditFileAtLine(path, lineNumber)
}

func (self *DiffExploreController) openFile() error {
	self.context.GetMutex().Lock()
	defer self.context.GetMutex().Unlock()

	path := self.selectedFilePath()
	if path == "" {
		return nil
	}
	return self.c.Helpers().Files.OpenFile(path)
}

func (self *DiffExploreController) escape() error {
	self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
	return nil
}

func (self *DiffExploreController) GetOnFocus() func(types.OnFocusOpts) {
	return func(opts types.OnFocusOpts) {
		self.c.Views().DiffExplore.Wrap = self.c.UserConfig().Gui.WrapLinesInStagingView
		self.prevExtrasVisible = self.c.State().GetShowExtrasWindow()
		self.c.State().SetShowExtrasWindow(false)
		self.c.Helpers().DiffExplore.RefreshDiffExplorePanel(opts)
	}
}

func (self *DiffExploreController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(opts types.OnFocusLostOpts) {
		self.context.SetState(nil)
		if opts.NewContextKey != self.context.GetKey() {
			self.c.Views().DiffExplore.Wrap = true
			self.c.State().SetShowExtrasWindow(self.prevExtrasVisible)
		}
	}
}
