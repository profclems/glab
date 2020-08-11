package commands

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

var pipelineViewCmd = &cobra.Command{
	Use:   "view <pipeline-id> [flags]",
	Short: `View a single pipeline`,
	Example: heredoc.Doc(`
	$ glab pipeline view 177883
	`),
	Long: ``,
	Run:  viewPipelines,
}

var pipelineDetails *gitlab.Pipeline
var pipelineJobDetails []*gitlab.Job
var mainView *gocui.View

func getPipelineJobs(pid int) []*gitlab.Job  {
	gitlabClient, repo := git.InitGitlabClient()
	l := &gitlab.ListJobsOptions{}
	pipeJobs, _, err := gitlabClient.Jobs.ListPipelineJobs(repo, pid, l)
	if err != nil {
		er(err)
	}
	return pipeJobs
}

func getPipelineJobLog(jobID int) io.Reader {
	gitlabClient, repo := git.InitGitlabClient()
	pipeJobs, _, err := gitlabClient.Jobs.GetTraceFile(repo, jobID)
	if err != nil {
		er(err)
	}
	return pipeJobs
}

func viewPipelines(cmd *cobra.Command, args []string) {
	if len(args) > 1 || len(args) == 0 {
		cmdErr(cmd, args)
		return
	}
	pid := manip.StringToInt(args[0])
	gitlabClient, repo := git.InitGitlabClient()
	fmt.Println("Finding pipeline...", pid)
	pipes, _, err := gitlabClient.Pipelines.GetPipeline(repo, pid)
	pipelineDetails = pipes
	if err != nil {
		er(err)
	}
	fmt.Println("Getting Pipeline Job...")
	pipelineJobDetails = getPipelineJobs(pid)
	fmt.Println("Setting up view...")
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.InputEsc = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("side")
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func showLoading(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Loading job log for", l)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	updatePipelineLog(mainView, 670715799)
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}
	/*
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	 */
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, showLoading); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlW, gocui.ModNone, saveVisualMain); err != nil {
		return err
	}
	return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	p := make([]byte, 5)
	v.Rewind()
	for {
		n, err := v.Read(p)
		if n > 0 {
			if _, err := f.Write(p[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveVisualMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	vb := v.ViewBuffer()
	if _, err := io.Copy(f, strings.NewReader(vb)); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	var firstJobID int
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for i, pipelineJob := range pipelineJobDetails {
			if i == 0 {
				firstJobID = pipelineJob.ID
			}
			fmt.Fprintf(v, "%s (%s)\n", pipelineJob.Name, pipelineJob.Status)
		}
	}
	return displayPipelineJobLog(g, firstJobID, maxX, maxY)
}

func displayPipelineJobLog(g *gocui.Gui, jid int, maxX, maxY int) error {
	if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		updatePipelineLog(v, jid)
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		v.Overwrite = true
		if _, err := g.SetCurrentView("side"); err != nil {
			return err
		}
		mainView = v
	}
	return nil
}

func updatePipelineLog(v *gocui.View, jid int)  {
	var str string
	if b, err := ioutil.ReadAll(getPipelineJobLog(jid)); err == nil {
		str = string(b)
	}
	str, _ = manip.RenderMarkdown(str)
	fmt.Fprintln(v, str)
	//updatePipelineLog(v, jid)
}

func init() {
	pipelineCmd.AddCommand(pipelineViewCmd)
}