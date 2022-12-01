# Releasing

This project uses Fyne's packaging to build its distributables. To install the necessary installing for packaging, run, from the root of this project:

```
make release-deps
```

Once that step has completed successfully, to build release artifacts for this project, run, from the root of this project:

```
make release
```

It will build targeted releases in the `dist/` folder beneath the root of this project.