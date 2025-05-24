# Expand Directory Dialog Size in Path Tab

## Summary
- Improved the directory selection dialog in the Path tab to open with a larger, more usable size
- Dialog now opens at 80% of the window size with a minimum of 800x600 pixels
- Provides better user experience when browsing for directories to add to PATH

## Changes
- Modified `shellconfig_gui.go` to use `dialog.NewFolderOpen` instead of `dialog.ShowFolderOpen`
- Added logic to calculate dialog size based on window dimensions
- Set minimum dimensions to ensure usability on smaller screens

## Testing
1. Open the Swiss Linux Knife application
2. Navigate to the Path tab
3. Click "Add Directory..."
4. Verify the dialog opens with a larger size (approximately double the default)
5. Verify you can still browse and select directories normally
6. Test on different window sizes to ensure the 80% scaling works correctly

## Screenshots
[Add before/after screenshots if available]

## Type of Change
- [ ] Bug fix
- [x] New feature
- [ ] Breaking change
- [ ] Documentation update

## Checklist
- [x] Code follows the style guidelines of this project
- [x] Self-review of code completed
- [x] Comments added where necessary
- [x] Changes generate no new warnings
- [x] Tested on local environment