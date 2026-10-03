## v0.8.0 (Sat, 03 Oct 2026 12:40:09 UTC)
- chore: update `github.com/ctx42/testing` to v0.56.0 and `github.com/ctx42/testkit` to v0.15.0.
- fix(ringtest): keep filesystem on Tester.Ring.
- docs(ring): document the ring package.
- docs(ringtest): document the ringtest package.
- docs(ring): describe metadata sentinels as caller errors.
- docs(ring): say which Clone fields are shared.
- fix(ring): skip a nil option passed to New.
- fix(ring): tolerate empty args, a nil clock, and a zero Ring.
- test(ring): widen the NowUTC window to one second.
- test(ringtest): expect dry-buffer writes on New.
- test(ring): order WithClock ahead of WithMeta.
- test(ring): order Test_EnvUnset ahead of its table test.
- test(ring): call IO getters from the When step.
- test(ring): separate IO test subjects with blank lines.
- test(ring): separate ring test subjects with blank lines.
- test(ring): separate env test subjects with blank lines.
- test(ringtest): separate tester subjects with blank lines.
- test(ring): name the NewIO result have.
- test(ring): name the New result have.
- test(ringtest): name When results have.
- test(ring): prepare When arguments.
- test(ring): prepare env When arguments.
- test(ringtest): prepare When arguments.
- test(ring): discard example print results.
- test(ringtest): align dry-buffer failure text.
- refactor(ring): assert Ring implements Environ.
- test(ring): assert Clone copies env and shares streams.
- docs(ringtest): say New stores nil arguments.
- docs: correct godoc grammar.
- docs(ring): cross-reference IO from NewIO.
- docs(ring): write Streamer method comments as sentences.
- test(ring): drop the redundant loop variable copy.
- fix(ring)!: return an empty slice from EnvAll.
- docs: give the README a runnable quickstart.
- fix(ring): allocate the map in EnvSet when it is nil.
- fix(ring)!: return an empty slice from SetFrom of a nil env.
- test(ring): assert EnvClone copies the entries.
- test(ring): reject a nil slice from an empty EnvUnset.
- test(ring): pin EnvOrOs of an empty slice.
- docs(ring): say a nil WithEnv is empty.
- docs(ring): warn that the zero Ring is not safe.
- docs(ring): mention the filesystem in the overview.
- docs(ring): write Ring comments as sentences.
- docs(ring): say where Ring.Name comes from.
- test(ring): separate the SetArgs assertions.
- docs(ring): correct the EnvSet interface wording.
- docs(ring): say EnvSet and EnvUnset return a new slice.
- docs(ring): write the Env compile-time check as a sentence.
- docs(ring): explain the Env method prefixes.
- docs(ring): write IO comments as sentences.
- docs(ring): say the IO setters store a nil stream.
- docs(ring): explain the IOClone name.
- docs(ringtest): describe tester streams on Ring.
- docs(ringtest): say wet buffers apply to the next ring.
- docs(ringtest): write Tester field comments as sentences.
- test(ringtest): expect nil arguments from New.
- test(ringtest): unwrap the filesystem readback.
- test(ringtest): separate assertion subjects.
- test(ringtest): name dry wet-buffer failures as errors.
- refactor(ring): drop redundant embed selectors.

## v0.7.1 (Fri, 03 Jul 2026 15:31:28 UTC)
- chore: update dependencies.

## v0.7.0 (Sun, 07 Jun 2026 20:51:02 UTC)
- fix: correct godoc and EnvSet/EnvUnset slice aliasing.
- chore: update deps, tooling config.
- test: sort EnvSet/EnvUnset results before comparing.
- docs: add testable example and fix README code blocks.
- docs: add examples for metadata, clone, and clock injection.

## v0.6.1 (Sat, 09 May 2026 19:56:20 UTC)
- chore: update github.com/ctx42/testing to v0.48.0.

## v0.6.0 (Fri, 01 May 2026 20:08:11 UTC)
- chore: Update to Go 1.26 and update dependencies.

## v0.5.3 (Thu, 30 Apr 2026 16:32:36 UTC)
- chore: Update dependencies.

## v0.5.2 (Mon, 15 Sep 2025 11:43:46 UTC)
- chore: Update dependencies.

## v0.5.1 (Thu, 07 Aug 2025 10:43:33 UTC)
- chore: update dependencies.

## v0.5.0 (Thu, 07 Aug 2025 10:42:47 UTC)
- feat: add fs.FS support to the Ring.

## v0.4.5 (Fri, 18 Jul 2025 20:30:14 UTC)
- chore: update dependencies.

## v0.4.4 (Mon, 07 Jul 2025 10:29:38 UTC)
- Update dependencies.

## v0.4.3 (Sat, 21 Jun 2025 07:30:15 UTC)
- Update dependencies.

## v0.4.2 (Fri, 20 Jun 2025 13:31:45 UTC)
- Update dependencies.

## v0.4.1 (Thu, 12 Jun 2025 07:47:13 UTC)
- Update dependencies.

## v0.4.0 (Mon, 09 Jun 2025 07:22:50 UTC)
- Add Clone methods for Ring, IO, and Env structures. Add a convenience method `Env.EnvSetWith(src []string)`.

## v0.3.0 (Sun, 08 Jun 2025 17:36:34 UTC)
- Update dependencies and add Ring.SetArgs method.

## v0.2.0 (Fri, 06 Jun 2025 20:27:20 UTC)
- Add Ring metadata access methods and update documentation.

## v0.1.0 (Fri, 06 Jun 2025 13:50:45 UTC)
- Add source Excalidraw file.

## v0.0.0 (Fri, 06 Jun 2025 13:42:53 UTC)
- Initial commit.
- Add LICENSE.md file.
- Implement filesystem abstraction.
- Refactor / simplify ring metadata.
- Write documentation, add GitHub Actions, code cleanup.
- Refactor environment handling, add tests refine comments and naming for better readability and consistency.
- Update to ctx42/testing to v0.9.0.
- Update dependencies.
- Update dependencies.

