# Features Brainstorm Doc

This doc includes possible features I'm considering adding and also a list of logical items that I am currently implementing. Part of this is for tracking, but I also thought it would be cool to share what I'm working on bringing to the tool. Suggestions are welcomed!

## Possible Features:

- Allow users to specify a config file that contains necessary options for the checks. Provides more flexibility.
- Allow users to specify a regex pattern to ignore files or directories that match the pattern.
- CodeAudit should automatically ignore files in certain directories like node_modules and hidden directories.
- Should be able to process multiple file types at once.
- Add a config flag for allowing the user to specify if the report should be saved in cwd or in the downloads folder. Will need to update logic to allow for this. Currently only saves in downloads folder.
- Currently you get the file location and line number where a check failed. Also share the position in a line of code where the check failed for quicker fixing and debugging for user.
- Use a llm model like ChatGpt or Gemini to give better name suggestions for variables.
- Turn CLI tool into a downloaded tool or one that can be installed via a package manager like homebrew

## Currently Working On:

- Creating more in depth mock files in the mock_directory for testing purposes.
- Add unit tests for all checks and utils.
- Instead of using a bunch of defaults, require the user to provide atleast the file type, the check type and the root directory. Use defaults for the rest if needed or missing.
- Cleanup the main function and move some of the logic to a separate functions
- Much logic in checks function is repeated. Refactor to remove the repetition.
<br>

_\* I was going to save this to an internal document, but thought it would be cool to share what I am working on._