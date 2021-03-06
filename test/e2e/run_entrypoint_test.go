package integration

import (
	"os"

	. "github.com/containers/libpod/v2/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Podman run entrypoint", func() {
	var (
		tempdir    string
		err        error
		podmanTest *PodmanTestIntegration
	)

	BeforeEach(func() {
		tempdir, err = CreateTempDirInTempDir()
		if err != nil {
			os.Exit(1)
		}
		podmanTest = PodmanTestCreate(tempdir)
		podmanTest.Setup()
		podmanTest.SeedImages()
	})

	AfterEach(func() {
		podmanTest.Cleanup()
		f := CurrentGinkgoTestDescription()
		processTestResult(f)

	})

	It("podman run no command, entrypoint, or cmd", func() {
		dockerfile := `FROM docker.io/library/alpine:latest
ENTRYPOINT []
CMD []
`
		podmanTest.BuildImage(dockerfile, "foobar.com/entrypoint:latest", "false")
		session := podmanTest.Podman([]string{"run", "foobar.com/entrypoint:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(125))
	})

	It("podman run entrypoint", func() {
		SkipIfRemote()
		dockerfile := `FROM docker.io/library/alpine:latest
ENTRYPOINT ["grep", "Alpine", "/etc/os-release"]
`
		podmanTest.BuildImage(dockerfile, "foobar.com/entrypoint:latest", "false")
		session := podmanTest.Podman([]string{"run", "foobar.com/entrypoint:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(Equal(2))
	})

	It("podman run entrypoint with cmd", func() {
		SkipIfRemote()
		dockerfile := `FROM docker.io/library/alpine:latest
CMD [ "-v"]
ENTRYPOINT ["grep", "Alpine", "/etc/os-release"]
`
		podmanTest.BuildImage(dockerfile, "foobar.com/entrypoint:latest", "false")
		session := podmanTest.Podman([]string{"run", "foobar.com/entrypoint:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(Equal(4))
	})

	It("podman run entrypoint with user cmd overrides image cmd", func() {
		SkipIfRemote()
		dockerfile := `FROM docker.io/library/alpine:latest
CMD [ "-v"]
ENTRYPOINT ["grep", "Alpine", "/etc/os-release"]
`
		podmanTest.BuildImage(dockerfile, "foobar.com/entrypoint:latest", "false")
		session := podmanTest.Podman([]string{"run", "foobar.com/entrypoint:latest", "-i"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(Equal(5))
	})

	It("podman run entrypoint with user cmd no image cmd", func() {
		SkipIfRemote()
		dockerfile := `FROM docker.io/library/alpine:latest
ENTRYPOINT ["grep", "Alpine", "/etc/os-release"]
`
		podmanTest.BuildImage(dockerfile, "foobar.com/entrypoint:latest", "false")
		session := podmanTest.Podman([]string{"run", "foobar.com/entrypoint:latest", "-i"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(Equal(5))
	})

	It("podman run user entrypoint overrides image entrypoint and image cmd", func() {
		SkipIfRemote()
		dockerfile := `FROM docker.io/library/alpine:latest
CMD ["-i"]
ENTRYPOINT ["grep", "Alpine", "/etc/os-release"]
`
		podmanTest.BuildImage(dockerfile, "foobar.com/entrypoint:latest", "false")
		session := podmanTest.Podman([]string{"run", "--entrypoint=uname", "foobar.com/entrypoint:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(session.LineInOuputStartsWith("Linux")).To(BeTrue())

		session = podmanTest.Podman([]string{"run", "--entrypoint", "", "foobar.com/entrypoint:latest", "uname"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(session.LineInOuputStartsWith("Linux")).To(BeTrue())
	})

	It("podman run user entrypoint with command overrides image entrypoint and image cmd", func() {
		SkipIfRemote()
		dockerfile := `FROM docker.io/library/alpine:latest
CMD ["-i"]
ENTRYPOINT ["grep", "Alpine", "/etc/os-release"]
`
		podmanTest.BuildImage(dockerfile, "foobar.com/entrypoint:latest", "false")
		session := podmanTest.Podman([]string{"run", "--entrypoint=uname", "foobar.com/entrypoint:latest", "-r"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(session.LineInOuputStartsWith("Linux")).To(BeFalse())
	})
})
