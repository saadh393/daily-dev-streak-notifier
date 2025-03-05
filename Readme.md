# Daily Dev Streak Notifier

**Stay Consistent with Your Daily Development Learning**

This command-line application is designed to help you maintain a consistent streak of daily development reading and learning. It provides a visual reminder in your terminal to encourage you to engage with development-related content every day.

**Key Features:**

*   **Daily Reminder:**  Displays a message in your terminal to remind you about your daily dev learning streak.
*   **Simple CLI:** Easy to use with straightforward command-line interactions.
*   **Customizable:** (Future Feature) While not mentioned in current file. There is scope of Customization of messages, and streak data.
* **Lightweight:** small footprint, no background processes.

**How it Works:**

The `dailydev` application is a self-contained executable. When run, it displays a message in your terminal. The intent is for you to add this command to your shell configuration file (e.g., `.bashrc`, `.zshrc`) so that it runs every time you open a new terminal window. This creates a visual cue to keep your development learning streak top-of-mind.

**Getting Started:**

**1. Installation:**

   *   **Clone the Repository:**
     ```bash
     git clone <repository-url> # Replace with your actual repository URL
     cd daily-dev
     ```
   * **Build from Source:**
        If you want to build from the source code, use the following command:
        ```bash
        go build -o dailydev
        ```
        This will create an executable file named `dailydev` in your project directory.
    * **Executable:**
       Make `dailydev` executable if it isn't already:
       ```bash
       chmod +x ./dailydev
       ```

**2. Running the Application:**

   *   **Manual Run:**
       You can run the application directly from your terminal:
       ```bash
       ./dailydev
       ```

   * **Automatic Run at Terminal Start**
   *  To get the reminder every time you open your terminal, you need to add it to your shell's startup file.
        * **Bash (e.g., Linux, macOS with Bash):**
                ```bash
                echo "./path/to/dailydev" >> ~/.bashrc
                source ~/.bashrc
                ```
            - replace `./path/to/dailydev` with real path of your project.
        * **Zsh (e.g., macOS with Zsh):**
                ```bash
                echo "./path/to/dailydev" >> ~/.zshrc
                source ~/.zshrc
                ```
            - replace `./path/to/dailydev` with real path of your project.
    * **Explanation:**
        * `echo "./path/to/dailydev"`: This adds the command to run `dailydev` to the end of your shell startup file.
        * `>>`: This appends the command to the file.
        * `~/.bashrc` or `~/.zshrc`: These are your shell configuration files.
        * `source ~/.bashrc` or `source ~/.zshrc`: This reloads your shell configuration so the changes take effect immediately.

**Building the Project (For Developers):**

1.  **Prerequisites:**
    *   Go programming language installed on your system (version 1.18 or later recommended).
2. **Build**
    *   Navigate to the project directory in your terminal.
    *   Run the following command:
    ```bash
    go build -o dailydev
    ```
    This will compile the Go code and create an executable file named `dailydev` in the same directory.

**Future Development:**

*   **Streak Tracking:** Implement a system to track the user's actual reading streak and display it.
*   **Customizable Messages:** Allow users to configure the reminder message.
*   **Data Storage:** Use a file or database to persist streak data.
* **Configuration File**: Allow users to customize the application behaviour using the file.
* **Interactive CLI**: Let the CLI allow user to do some activity, like reset the streak.

**Contributing:**

Contributions are welcome! If you have any ideas for improving this project, please feel free to open an issue or submit a pull request.

**License:**

[Add your License here, e.g., MIT License]
