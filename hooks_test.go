package fiber

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3/utils"
	"github.com/valyala/bytebufferpool"
)

var testSimpleHandler = func(c *Ctx) error {
	return c.SendString("simple")
}

func Test_Hook_OnRoute(t *testing.T) {
	t.Parallel()

	app := New()

	app.Hooks().OnRoute(func(r Route) error {
		utils.AssertEqual(t, "", r.Name)

		return nil
	})

	app.Get("/", testSimpleHandler).Name("x")

	subApp := New()
	subApp.Get("/test", testSimpleHandler)

	app.Mount("/sub", subApp)
}

func Test_Hook_OnName(t *testing.T) {
	t.Parallel()

	app := New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Hooks().OnName(func(r Route) error {
		buf.WriteString(r.Name)

		return nil
	})

	app.Get("/", testSimpleHandler).Name("index")

	subApp := New()
	subApp.Get("/test", testSimpleHandler)
	subApp.Get("/test2", testSimpleHandler)

	app.Mount("/sub", subApp)

	utils.AssertEqual(t, "index", buf.String())
}

func Test_Hook_OnName_Error(t *testing.T) {
	t.Parallel()

	app := New()
	defer func() {
		if err := recover(); err != nil {
			utils.AssertEqual(t, "unknown error", fmt.Sprintf("%v", err))
		}
	}()

	app.Hooks().OnName(func(r Route) error {
		return errors.New("unknown error")
	})

	app.Get("/", testSimpleHandler).Name("index")
}

func Test_Hook_OnGroup(t *testing.T) {
	t.Parallel()

	app := New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Hooks().OnGroup(func(g Group) error {
		buf.WriteString(g.Prefix)

		return nil
	})

	grp := app.Group("/x").Name("x.")
	grp.Group("/a")

	utils.AssertEqual(t, "/x/x/a", buf.String())
}

func Test_Hook_OnGroupName(t *testing.T) {
	t.Parallel()

	app := New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Hooks().OnGroupName(func(g Group) error {
		buf.WriteString(g.name)

		return nil
	})

	grp := app.Group("/x").Name("x.")
	grp.Get("/test", testSimpleHandler)
	grp.Get("/test2", testSimpleHandler)

	utils.AssertEqual(t, "x.", buf.String())
}

func Test_Hook_OnGroupName_Error(t *testing.T) {
	t.Parallel()

	app := New()
	defer func() {
		if err := recover(); err != nil {
			utils.AssertEqual(t, "unknown error", fmt.Sprintf("%v", err))
		}
	}()

	app.Hooks().OnGroupName(func(g Group) error {
		return errors.New("unknown error")
	})

	grp := app.Group("/x").Name("x.")
	grp.Get("/test", testSimpleHandler)
}

func Test_Hook_OnShutdown(t *testing.T) {
	t.Parallel()

	app := New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Hooks().OnShutdown(func() error {
		buf.WriteString("shutdowning")

		return nil
	})

	utils.AssertEqual(t, nil, app.Shutdown())
	utils.AssertEqual(t, "shutdowning", buf.String())
}

func Test_Hook_OnListen(t *testing.T) {
	t.Parallel()

	app := New()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	app.Hooks().OnListen(func() error {
		buf.WriteString("ready")

		return nil
	})

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()
	utils.AssertEqual(t, nil, app.Start(":9000", StartConfig{DisableStartupMessage: true}))

	utils.AssertEqual(t, "ready", buf.String())
}
