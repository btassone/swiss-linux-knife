# Add PATH Order Numbers and Functional Move Buttons

## Summary
- Added visual order numbers to PATH entries showing search priority
- Implemented functional Move Up/Down buttons for reordering paths
- Improved user understanding of PATH search order

## Changes
- Added order labels (1., 2., 3., etc.) in bold next to each path
- Implemented table selection tracking
- Made Move Up/Down buttons functional
- Selected items stay selected after moving
- Path order is preserved when saving

## Benefits
- Users can now see the exact search order of their PATH
- Easy reordering with visual feedback
- Better understanding of how PATH resolution works
- Maintains selection for continuous reordering

## Testing
1. Open the Path tab
2. Verify each path shows its order number
3. Click on a path entry to select it
4. Use Move Up/Down buttons to reorder
5. Verify the selection follows the moved item
6. Verify order numbers update correctly
7. Save and reload to ensure order is preserved

## Type of Change
- [ ] Bug fix
- [x] New feature
- [ ] Breaking change
- [ ] Documentation update