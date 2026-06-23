# Pass TUI Tree View - Implementation Tasks

## Feature Number: 006
## Specification: [spec.md](./spec.md)

## Overview
This document outlines the implementation tasks for adding tree view to the pass TUI fuzzy finder.

## Task List

### Phase 1: Test-Driven Development ✅ IN PROGRESS

- [ ] **TASK-001**: Create test file `cmd/tui/tree_test.go`
  - [ ] Test tree to list conversion
  - [ ] Test tree rendering in TUI context
  - [ ] Test filtering with tree view
  - [ ] Test selection with tree view
  - **Estimate**: 2 hours
  - **Priority**: High
  - **Status**: Not Started

- [ ] **TASK-002**: Add tests to existing `cmd/tui/tui_test.go`
  - [ ] Test tree view in different modes
  - [ ] Verify flat view is not used
  - **Estimate**: 1 hour
  - **Priority**: High
  - **Status**: Not Started

### Phase 2: Core Implementation

- [ ] **TASK-003**: Create TreeItem struct in `cmd/tui/models.go`
  - [ ] Implement `list.Item` interface
  - [ ] Store path, display name, indentation
  - [ ] Support match highlighting
  - **Estimate**: 1 hour
  - **Priority**: High
  - **Status**: Not Started

- [ ] **TASK-004**: Add tree construction to TUI model
  - [ ] Import tree package
  - [ ] Build tree from passwords in `NewModel()`
  - [ ] Store tree root in Model struct
  - **Estimate**: 1 hour
  - **Priority**: High
  - **Status**: Not Started

- [ ] **TASK-005**: Implement tree to list conversion
  - [ ] Create `flattenTreeToListItems()` function
  - [ ] Handle indentation and connectors
  - [ ] Preserve path for filtering
  - **Estimate**: 2 hours
  - **Priority**: High
  - **Status**: Not Started

- [ ] **TASK-006**: Update Model struct
  - [ ] Add `treeRoot` field
  - [ ] Add `flatView` field (for future toggle)
  - [ ] Update `NewModel()` to build tree
  - **Estimate**: 1 hour
  - **Priority**: High
  - **Status**: Not Started

- [ ] **TASK-007**: Modify filtering for tree view
  - [ ] Filter on full path, not display name
  - [ ] Preserve match highlighting
  - [ ] Update `filterList()` function
  - **Estimate**: 2 hours
  - **Priority**: High
  - **Status**: Not Started

### Phase 3: Integration

- [ ] **TASK-008**: Integrate tree view into all modes
  - [ ] Verify show mode works with tree
  - [ ] Verify clip mode works with tree
  - [ ] Verify rm mode works with tree
  - [ ] Verify edit mode works with tree
  - **Estimate**: 1 hour
  - **Priority**: Medium
  - **Status**: Not Started

- [ ] **TASK-009**: Ensure backward compatibility
  - [ ] Verify all existing tests pass
  - [ ] Test CLI commands still work
  - [ ] Test `pass find` still uses flat view
  - **Estimate**: 1 hour
  - **Priority**: High
  - **Status**: Not Started

### Phase 4: Documentation

- [ ] **TASK-010**: Update TUI spec document
  - [ ] Document tree view feature
  - [ ] Update architecture diagram
  - **Estimate**: 1 hour
  - **Priority**: Medium
  - **Status**: Not Started

- [ ] **TASK-011**: Update README if needed
  - [ ] Add tree view description
  - [ ] Update screenshots if available
  - **Estimate**: 30 minutes
  - **Priority**: Low
  - **Status**: Not Started

### Phase 5: Testing & Validation

- [ ] **TASK-012**: Manual testing
  - [ ] Test with various directory structures
  - [ ] Test fuzzy matching with tree view
  - [ ] Test all TUI modes
  - [ ] Test edge cases (empty store, single password, deep nesting)
  - **Estimate**: 2 hours
  - **Priority**: High
  - **Status**: Not Started

- [ ] **TASK-013**: Performance testing
  - [ ] Verify tree construction doesn't slow down TUI
  - [ ] Test with large password stores (>1000 entries)
  - **Estimate**: 1 hour
  - **Priority**: Medium
  - **Status**: Not Started

## Implementation Order

Recommended order based on dependencies:

1. **TASK-001** - Create tree_test.go
2. **TASK-002** - Add tests to tui_test.go
3. **TASK-003** - Create TreeItem struct
4. **TASK-004** - Add tree construction to model
5. **TASK-005** - Implement tree to list conversion
6. **TASK-006** - Update Model struct
7. **TASK-007** - Modify filtering for tree view
8. **TASK-008** - Integrate tree view into all modes
9. **TASK-012** - Manual testing
10. **TASK-009** - Ensure backward compatibility
11. **TASK-013** - Performance testing
12. **TASK-010** - Update documentation
13. **TASK-011** - Update README

## File Changes Summary

| File | Change Type | Description |
|------|-------------|-------------|
| `cmd/tui/tree_test.go` | NEW | Tests for tree view in TUI |
| `cmd/tui/models.go` | MODIFY | Add TreeItem, tree construction, tree to list conversion |
| `cmd/tui/tui_test.go` | MODIFY | Add tree view tests |
| `specs/pass-tui-spec.md` | MODIFY | Document tree view feature |
| `README` | MODIFY (optional) | Update with tree view info |

## Dependencies

- **Feature 004**: Pass Tree View for Find Command (already implemented)
  - Provides `cmd/tree/tree.go` package with tree rendering
  - Must be complete before starting this feature

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Tree construction slows down TUI | Low | Medium | Optimize tree building, test with large stores |
| Breaking existing TUI functionality | Medium | High | Extensive testing, maintain backward compatibility |
| Filtering doesn't work with tree view | Medium | High | Test filtering thoroughly, preserve path-based matching |
| Tree rendering issues | Medium | Medium | Test various directory structures |

## Success Metrics

- All tests pass (existing + new)
- Tree view works in all TUI modes
- `pass find` maintains flat output
- No performance degradation
- No breaking changes

## Checklist for Completion

- [ ] All tasks completed
- [ ] All tests pass
- [ ] Manual testing successful
- [ ] Documentation updated
- [ ] Code reviewed
- [ ] Merged to main branch

## Notes

- Follow spec-driven development: write tests first, then implementation
- Use existing tree package from feature 004
- Maintain backward compatibility
- Keep changes focused and minimal
- Update documentation as you go

---

*Feature Number: 006*
*Last Updated: 2026-06-22*
*Status: Ready for Implementation*