#!/bin/bash

# Initialize an empty array
html_files=()

# Use a while loop to read the output of find into html_files array
while IFS= read -r file; do
    html_files+=("$file")
done < <(find . -maxdepth 1 -type f -name "*.html")

# Check if there are HTML files in the current directory
if [ ${#html_files[@]} -eq 0 ]; then
  echo "No HTML files found in the current directory."
  exit 1
fi

echo "Found the following HTML files:"
for i in "${!html_files[@]}"; do
  echo "$((i+1))) ${html_files[$i]}"
done

# Prompt the user to select a file
echo "Enter the number of the HTML file you want to process:"
read file_number

# Validate input
if [[ ! $file_number =~ ^[0-9]+$ ]] || [ $file_number -lt 1 ] || [ $file_number -gt ${#html_files[@]} ]; then
  echo "Invalid selection."
  exit 1
fi

# Adjust index to match array indexing starting at 0
let file_number=file_number-1

selected_file=${html_files[$file_number]}

# Prompt the user for a new file name
echo "Enter a new file name (including .html extension) for the processed file:"
read new_file_name

# Validate new file name is not empty
if [ -z "$new_file_name" ]; then
  echo "The file name cannot be empty."
  exit 1
fi

# Use sed to replace text outside HTML tags and save to the new file
sed -E 's/>([^<]+)</><!-- TEXT HERE -->\n</g' "$selected_file" > "$new_file_name"

echo "Processing complete. Modified file saved as $new_file_name"