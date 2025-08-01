package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	fmt.Println("🚀 GUI starting...")

	myApp := app.New()
	myWindow := myApp.NewWindow("Task Manager GUI")

	filename := "stored_tasks.txt"

	background := canvas.NewImageFromFile("background.png")
	background.FillMode = canvas.ImageFillStretch

	titleLabel := canvas.NewText("Task Manager GUI", color.Black)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 22

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title")

	headerEntry := widget.NewEntry()
	headerEntry.SetPlaceHolder("Header")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Description")

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetPlaceHolder("Output will appear here")

	watermark := canvas.NewText("Made by Navami Borkar", color.RGBA{R: 120, G: 120, B: 120, A: 255})
	watermark.TextStyle = fyne.TextStyle{Bold: true}
	watermark.TextSize = 14

	addBtn := widget.NewButton("Add Task", func() {
		title := strings.TrimSpace(titleEntry.Text)
		header := strings.TrimSpace(headerEntry.Text)
		description := strings.TrimSpace(descEntry.Text)

		if title == "" || header == "" || description == "" {
			output.SetText("Please fill in all fields before adding a task.")
			return
		}

		tasks, err := LoadTasks(filename)
		if err != nil {
			output.SetText("Error loading tasks: " + err.Error())
			return
		}

		for _, task := range tasks {
			if task.Title == title {
				output.SetText("Task with this title already exists.")
				return
			}
		}

		newTask := Task{
			Title:       title,
			Header:      header,
			Description: description,
		}
		tasks = append(tasks, newTask)

		if err := SaveTasks(filename, tasks); err != nil {
			output.SetText("Error saving task: " + err.Error())
		} else {
			output.SetText(fmt.Sprintf("Task added successfully: %s", title))
			titleEntry.SetText("")
			headerEntry.SetText("")
			descEntry.SetText("")
		}
	})

	showBtn := widget.NewButton("Show All Tasks", func() {
		tasks, err := LoadTasks(filename)
		if err != nil {
			output.SetText("Error: " + err.Error())
			return
		}
		if len(tasks) == 0 {
			output.SetText("No tasks found.")
			return
		}

		var builder strings.Builder
		for idx, task := range tasks {
			builder.WriteString(fmt.Sprintf("Index: %d\nTitle: %s\nHeader: %s\nDescription: %s\n---\n",
				idx, task.Title, task.Header, task.Description))
		}
		output.SetText(builder.String())
	})

	updateBtn := widget.NewButton("Update Task", func() {
		tasks, err := LoadTasks(filename)
		if err != nil {
			output.SetText("Error loading tasks: " + err.Error())
			return
		}

		title := strings.TrimSpace(titleEntry.Text)
		if title == "" {
			output.SetText("Please enter the title of the task to update.")
			return
		}

		updated := false
		for i, task := range tasks {
			if task.Title == title {
				tasks[i].Header = headerEntry.Text
				tasks[i].Description = descEntry.Text
				updated = true
				break
			}
		}

		if !updated {
			output.SetText("Task not found.")
			return
		}

		if err := SaveTasks(filename, tasks); err != nil {
			output.SetText("Error saving updated task: " + err.Error())
		} else {
			output.SetText("Task updated successfully.")
		}
	})

	deleteBtn := widget.NewButton("Delete Task", func() {
		tasks, err := LoadTasks(filename)
		if err != nil {
			output.SetText("Error loading tasks: " + err.Error())
			return
		}

		title := strings.TrimSpace(titleEntry.Text)
		if title == "" {
			output.SetText("Please enter the title of the task to delete.")
			return
		}

		newTasks := []Task{}
		deleted := false
		for _, task := range tasks {
			if task.Title == title {
				deleted = true
				continue // skip this task
			}
			newTasks = append(newTasks, task)
		}

		if !deleted {
			output.SetText("Task not found.")
			return
		}

		if err := SaveTasks(filename, newTasks); err != nil {
			output.SetText("Error deleting task: " + err.Error())
		} else {
			output.SetText("Task deleted successfully.")
		}
	})

	content := container.NewVBox(
		container.NewCenter(titleLabel),
		titleEntry,
		headerEntry,
		descEntry,
		addBtn,
		showBtn,
		updateBtn,
		deleteBtn,
		output,
		container.NewCenter(watermark),
	)

	myWindow.SetContent(container.NewMax(background, content))
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.ShowAndRun()
}
