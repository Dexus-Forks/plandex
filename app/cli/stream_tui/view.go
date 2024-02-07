package streamtui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var borderColor = lipgloss.Color("#444")
var helpTextColor = lipgloss.Color("#ddd")

func (m streamUIModel) View() string {
	if m.promptingMissingFile {
		return m.renderMissingFilePrompt()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderMainView(),
		m.renderProcessing(),
		m.renderBuild(),
		m.renderHelp(),
	)
}

func (m streamUIModel) renderMainView() string {
	return m.replyViewport.View()
}

func (m streamUIModel) renderHelp() string {
	style := lipgloss.NewStyle().Width(m.width).Foreground(lipgloss.Color(helpTextColor)).BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderForeground(lipgloss.Color(borderColor))

	return style.Render(" (s)top")
}

func (m streamUIModel) renderProcessing() string {
	if m.starting || m.processing {
		style := lipgloss.NewStyle().Width(m.width).BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderForeground(lipgloss.Color(borderColor))

		return style.Render(m.spinner.View())
	} else {
		return ""
	}
}

func (m streamUIModel) renderBuild() string {
	if !m.building {
		return ""
	}

	style := lipgloss.NewStyle().Width(m.width).BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderForeground(lipgloss.Color(borderColor))

	head := color.New(color.BgGreen, color.FgHiWhite, color.Bold).Sprint(" 🏗  ") + color.New(color.BgGreen, color.FgHiWhite).Sprint("Building plan ")

	filePaths := make([]string, 0, len(m.tokensByPath))
	for filePath := range m.tokensByPath {
		filePaths = append(filePaths, filePath)
	}

	sort.Strings(filePaths)

	var rows [][]string
	lineWidth := 0
	lineNum := 0

	for _, filePath := range filePaths {
		tokens := m.tokensByPath[filePath]
		finished := m.finishedByPath[filePath]
		block := fmt.Sprintf("📄 %s", filePath)

		if finished {
			block += " ✅"
		} else if tokens > 0 {
			block += fmt.Sprintf(" %d 🪙", tokens)
		}

		blockWidth := lipgloss.Width(block)

		if lineWidth+blockWidth > m.width {
			lineWidth = 0
			lineNum++
		}

		if len(rows) <= lineNum {
			rows = append(rows, []string{})
		}

		row := rows[lineNum]
		row = append(row, block)
		rows[lineNum] = row

		lineWidth += blockWidth
	}

	resRows := make([]string, len(rows)+1)

	resRows[0] = head
	for i, row := range rows {
		resRows[i+1] = strings.Join(row, " | ")
	}

	return style.Render(strings.Join(resRows, "\n"))
}

func (m streamUIModel) renderMissingFilePrompt() string {
	style := lipgloss.NewStyle().Padding(1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(borderColor)).Width(m.width - 2).Height(m.height - 2)

	prompt := "📄 " + color.New(color.Bold, color.FgHiYellow).Sprint(m.missingFilePath) + " isn't in context."

	prompt += "\n\n"

	desc := "This file exists in your project, but isn't loaded into context. Unless you load it into context or skip generating it, Plandex will fully overwrite the existing file rather than applying updates."

	words := strings.Split(desc, " ")
	for i, word := range words {
		words[i] = color.New(color.FgWhite).Sprint(word)
	}

	prompt += strings.Join(words, " ")

	prompt += "\n\n" + color.New(color.FgHiMagenta, color.Bold).Sprintln("🧐 What do you want to do?")

	for i, opt := range missingFileSelectOpts {
		if i == m.missingFileSelectedIdx {
			prompt += color.New(color.FgHiCyan, color.Bold).Sprint(" > " + opt)
		} else {
			prompt += "   " + opt
		}

		if opt == MissingFileLoadLabel {
			prompt += fmt.Sprintf(" | %d 🪙", m.missingFileTokens)
		}

		prompt += "\n"
	}

	return style.Render(prompt)
}
