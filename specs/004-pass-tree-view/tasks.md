# Pass Tree View Implementation Tasks

## Overview

This document tracks the implementation tasks for the tree view feature (Feature 004) in `pass find` command as specified in `specs/pass-tree-view-004/spec.md`.

## Task Breakdown

### Phase 1: Specification ✅
- [x] Create spec document (`specs/pass-tree-view/spec.md`)
- [x] Create tasks document (`specs/pass-tree-view/tasks.md`)
- [x] Review spec with stakeholders

### Phase 2: Tree Package Implementation ✅

#### 2.1 Create Tree Package Structure
- [x] Create `pass/cmd/tree/tree.go`
- [x] Create `pass/cmd/tree/tree_test.go`

#### 2.2 Implement Tree Node Structure ✅
- [x] Define `TreeNode` struct with Name, IsDir, Children fields
- [x] Implement `NewTreeNode()` constructor
- [x] Implement `AddChild()` with alphabetical sorting
- [x] Implement `FindOrCreateChild()` helper

#### 2.3 Implement Tree Rendering ✅
- [x] Implement `Render()` method with box-drawing characters
- [x] Handle prefix building for nested levels
- [x] Handle connector selection (├── vs └──)
- [x] Handle indentation (4 spaces per level)
- [x] Format directory names with trailing `/`

#### 2.4 Implement Tree Construction ✅
- [x] Implement `BuildTreeFromPaths()` function
- [x] Parse paths into components
- [x] Build tree structure from path list
- [x] Handle `.gpg` extension stripping
- [x] Mark directories vs files correctly

### Phase 3: Find Command Integration ✅

#### 3.1 Add Flags to Find Command
- [x] Add `--flat, -f` flag to `findCmd`
- [x] Add `--no-tree` flag to `findCmd`
- [x] Update flag variable declarations

#### 3.2 Modify Find Function ✅
- [x] Update `findPasswords()` signature to accept `flat bool` parameter
- [x] Add tree view rendering logic with `renderTreeNode()`
- [x] Keep flat view as fallback/option
- [x] Handle edge cases (empty results, single result)

#### 3.3 Update Root Command ✅
- [x] Ensure `addFindCmd()` is called in `Execute()`
- [x] Verify flag registration

### Phase 4: Documentation Updates ✅

#### 4.1 Update TUI Skill Documentation
- [x] Add tree view rendering section to `docs/tui-skill.md`
- [x] Add example of tree structure with Bubble Tea
- [x] Document box-drawing characters usage
- [x] Add styling examples for tree views

#### 4.2 Update Pass Documentation
- [ ] Update `README` (root) to mention tree view for find
- [ ] Document new flags in usage examples

### Phase 5: Testing

#### 5.1 Unit Tests ✅
- [x] Create unit tests for `tree/tree.go`
  - [x] Test `NewTreeNode()`
  - [x] Test `AddChild()` with sorting
  - [x] Test `FindOrCreateChild()`
  - [x] Test `Render()` with various structures
  - [x] Test `BuildTreeFromPaths()`
- [x] Update `cmd/find_test.go`
  - [x] Test tree view output
  - [x] Test flat view still works
  - [x] Test flag parsing
  - [x] Test edge cases

#### 5.2 Integration Tests
- [x] Add tree view test to `cmd/find_test.go`
  - [x] Test find with tree view
  - [x] Test find with flat view
  - [x] Test with nested directory structures

#### 5.3 Manual Testing
- [ ] Test on Windows
- [ ] Test on Linux/macOS
- [ ] Test with various directory structures
- [ ] Test with special characters in paths
- [ ] Test with empty results
- [ ] Test flag combinations

### Phase 6: Finalization

#### 6.1 Code Review
- [ ] Review tree package implementation
- [ ] Review find command changes
- [ ] Check error handling
- [ ] Verify no breaking changes to existing functionality

#### 6.2 Cleanup
- [ ] Remove debug code
- [ ] Fix linting issues
- [ ] Ensure consistent code style

#### 6.3 Commit
- [ ] Stage all changes
- [ ] Write commit message with gitmoji
- [ ] Commit to branch

## Priority Order

1. **High Priority (Must Have)**
   - Tree node structure and rendering
   - Tree construction from paths
   - Find command integration
   - Basic tests

2. **Medium Priority (Should Have)**
   - Flag support (--flat, --no-tree)
   - Documentation updates
   - Integration tests

3. **Low Priority (Nice to Have)**
   - Color/styling for tree view
   - TUI skill documentation updates

## Estimated Time

| Task | Estimate |
|------|----------|
| Spec and planning | 1 hour |
| Tree package implementation | 2 hours |
| Find command integration | 1 hour |
| Unit tests | 1 hour |
| Integration tests | 1 hour |
| Documentation | 1 hour |
| Finalization | 30 minutes |
| **Total** | **7.5 hours** |

## Dependencies

- Go 1.20+
- Existing pass tool structure
- No new external dependencies

## Blockers

None identified.

## Notes

- Follow existing code style and conventions in pass package
- Ensure backward compatibility (--flat flag preserves original behavior)
- Tree view should match the visual style shown in the Unix pass screenshot
- All box-drawing characters should render correctly in standard terminals
