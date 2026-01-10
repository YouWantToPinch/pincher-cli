// Package ui provides custom bubbletea models intended for the pincher-cli
package ui

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// TModelStructMenu is a bubbletea model that can be used to expose
// primitive struct fields to end users for input,
// as if they were elements of a menu.
type TModelStructMenu struct {
	fields []string // fields which can be edited; populated dynamically
	cursor int      // which field our cursor is pointing at
	// tracks state of field editing
	isEditingValue bool
	structType     reflect.Type
	structFields   map[string]any // field values
	QuitWithCancel bool           // can be used to communicate whether changes ought be saved
}

func InitialTModelStructMenu(structObj any, fieldList []string, asBlacklist bool) (TModelStructMenu, error) {
	// if fieldList is empty, all fields are exposed to users; otherwise, it is used as a whitelist.
	// if bool parameter 'asBlacklist' is 'true', the fieldList is used as a blacklist instead of a whitelist.
	t := reflect.TypeOf(structObj)
	v := reflect.ValueOf(structObj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	} else {
		return TModelStructMenu{}, errors.New("structObj should be a pointer to struct, so as to have addressable fields")
	}
	if t.Kind() != reflect.Struct {
		fmt.Println("ERROR: Not a struct. Check your input!")
		return TModelStructMenu{}, nil
	}
	newModel := TModelStructMenu{
		isEditingValue: false,
		structType:     t,
		structFields:   make(map[string]any),
		QuitWithCancel: false,
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if len(fieldList) != 0 {
			if asBlacklist {
				if slices.Contains(fieldList, field.Name) {
					continue
				}
			} else {
				if !(slices.Contains(fieldList, field.Name)) {
					continue
				}
			}
		}

		fieldVal := v.FieldByName(field.Name)
		if !fieldVal.CanSet() {
			fmt.Printf("Warning: Field '%s' left unexposed (cannot be set; unexported or not addressable).\n", field.Name)
			continue
		}

		if field.Type.Kind() == reflect.String {
			newModel.fields = append(newModel.fields, field.Name)
			newModel.structFields[field.Name] = fieldVal.Interface() // ""
		} else if field.Type.Kind() == reflect.Bool {
			newModel.fields = append(newModel.fields, field.Name)
			newModel.structFields[field.Name] = fieldVal.Interface() // false
		} else if field.Type.Kind() >= reflect.Int && field.Type.Kind() <= reflect.Int64 {
			newModel.fields = append(newModel.fields, field.Name)
			newModel.structFields[field.Name] = fieldVal.Interface() // 0
		}
	}

	if len(newModel.structFields) == 0 {
		return TModelStructMenu{}, fmt.Errorf("ERROR: No fields to expose to users in struct '%v'", newModel.structType.Name())
	}

	return newModel, nil
}

func (m TModelStructMenu) ParseStruct(obj any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("ERROR: expected a pointer to a struct, got %v", v.Kind())
	}
	v = v.Elem()

	for fieldName, newValue := range m.structFields {
		field := v.FieldByName(fieldName)

		if !field.IsValid() {
			fmt.Printf("Warning: Field '%s' not found in struct.\n", fieldName)
			continue
		}
		if !field.CanSet() {
			fmt.Printf("Warning: Field '%s' cannot be set (unexported or not addressable).\n", fieldName)
			continue
		}

		if field.Kind() >= reflect.Int && field.Kind() <= reflect.Int64 {
			if val, ok := newValue.(int); ok {
				field.SetInt(int64(val))
			} else {
				return fmt.Errorf("type mismatch for field '%s': expected int, got %T", fieldName, newValue)
			}
		} else if field.Kind() == reflect.Bool {
			if val, ok := newValue.(bool); ok {
				field.SetBool(val)
			} else if val, ok := newValue.(int); ok {
				boolVal := (val != 0)
				// fmt.Println(fmt.Sprintf("Bool digit value %d translated as: %t", val, boolVal))
				field.SetBool(boolVal)
			} else if val, ok := newValue.(string); ok {
				boolVal := (val != "f")
				// fmt.Println(fmt.Sprintf("Bool string value %s translated as: %t", val, boolVal))
				field.SetBool(boolVal)
			} else if !ok {
				fmt.Println("Error parsing digit as boolean value.")
			}
		} else if field.Kind() == reflect.String {
			if val, ok := newValue.(string); ok {
				field.SetString(val)
			}
		} else {
			fmt.Printf("Skipping field '%s': unsupported kind %s\n", fieldName, field.Kind())
		}
	}
	return nil
}

func (m TModelStructMenu) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m TModelStructMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// toggle edit mode on field if 'enter' key was pressed
		if msg.String() == "enter" {
			m.isEditingValue = !(m.isEditingValue)
		} else if msg.Type == tea.KeyBackspace {
			switch m.structFields[m.fields[m.cursor]].(type) {

			case string:
				stringVal := m.structFields[m.fields[m.cursor]].(string)
				if len(stringVal) > 0 {
					m.structFields[m.fields[m.cursor]] = stringVal[:len(stringVal)-1]
				}
			case int:
				if m.structFields[m.fields[m.cursor]].(int) != 0 {
					intSign := 1
					if m.structFields[m.fields[m.cursor]].(int) < 0 {
						intSign = -1
					}
					stringVal := strconv.Itoa(m.structFields[m.fields[m.cursor]].(int))
					var newVal string
					if intSign == 1 {
						newVal = stringVal[:len(stringVal)-1]
					} else {
						newVal = stringVal[1 : len(stringVal)-1]
					}
					if len(newVal) == 0 {
						m.structFields[m.fields[m.cursor]] = 0
					} else {
						convValue, err := strconv.Atoi(newVal)
						if err != nil {
							fmt.Printf("ERROR converting ascii to int: %v\n", err)
						} else {
							m.structFields[m.fields[m.cursor]] = (convValue * intSign)
						}
					}
				}
			}
		} else {
			if m.isEditingValue {
				switch m.structFields[m.fields[m.cursor]].(type) {
				case bool:
					switch msg.String() {
					case "t", "1":
						m.structFields[m.fields[m.cursor]] = true
					case "f", "0":
						m.structFields[m.fields[m.cursor]] = false
					case "right", "left":
						m.structFields[m.fields[m.cursor]] = !(m.structFields[m.fields[m.cursor]].(bool))
					default:
						m.structFields[m.fields[m.cursor]] = false
					}

				case string:
					m.structFields[m.fields[m.cursor]] = m.structFields[m.fields[m.cursor]].(string) + msg.String()
				case int:
					switch msg.String() {

					// The "right" and "l" keys increase the value
					case "right", "l":
						m.structFields[m.fields[m.cursor]] = m.structFields[m.fields[m.cursor]].(int) + 1

					// The "left" and "h" keys decrease the value
					case "left", "h":
						m.structFields[m.fields[m.cursor]] = m.structFields[m.fields[m.cursor]].(int) - 1

					case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
						if m.structFields[m.fields[m.cursor]] == 0 {
							convValue, err := strconv.Atoi(msg.String())
							if err != nil {
								fmt.Printf("ERROR: failed to convert ascii to int: %v\n", err)
							} else {
								m.structFields[m.fields[m.cursor]] = convValue
							}
						} else {
							intValue, err := strconv.Atoi(strconv.Itoa(m.structFields[m.fields[m.cursor]].(int)) + msg.String())
							if err != nil {
								fmt.Printf("ERROR: %v\n", err)
							}
							m.structFields[m.fields[m.cursor]] = intValue
						}
					}
				}
			} else {
				// Cool, what was the actual key pressed?
				switch msg.String() {

				case "s":
					return m, tea.Quit

				// These keys should exit the program.
				case "ctrl+c", "q":
					m.QuitWithCancel = true
					return m, tea.Quit

				// The "up" and "k" keys move the cursor up
				case "up", "k":
					if m.cursor > 0 {
						m.cursor--
					}

				// The "down" and "j" keys move the cursor down
				case "down", "j":
					if m.cursor < len(m.fields)-1 {
						m.cursor++
					}

				// Any numeric key sets the value for the item that
				// the cursor is pointing at.
				case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
					intValue, err := strconv.Atoi(msg.String())
					if err != nil {
						fmt.Printf("ERROR: %v\n", err)
					}
					m.structFields[m.fields[m.cursor]] = intValue
				}
			}
		}
	}

	// Return the updated TModelStructMenu to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m TModelStructMenu) View() string {
	// The header
	s := "Set values for the following:\n\n"

	// for formatting, get longest field name
	maxFieldName := 0
	for _, field := range m.fields {
		if len(field) > maxFieldName {
			maxFieldName = len(field)
		}
	}

	// Iterate over our fields
	for i, choice := range m.fields {

		// Is the cursor pointing at this choice?
		cursor := "  " // no cursor
		if m.cursor == i {
			if m.isEditingValue {
				cursor = ">>"
			} else {
				cursor = "> "
			}
		}

		// Is this choice numerated?
		var value string // string represenation of field value
		if _, ok := m.structFields[m.fields[i]]; ok {
			switch m.structFields[m.fields[i]].(type) {
			case string:
				if m.isEditingValue && m.cursor == i {
					iBeam := "|"
					value = m.structFields[m.fields[i]].(string) + iBeam
				} else {
					value = m.structFields[m.fields[i]].(string)
				}
			case bool:
				value = strconv.FormatBool(m.structFields[m.fields[i]].(bool))
			case int:
				value = strconv.Itoa(m.structFields[m.fields[i]].(int))
			}
		}

		// Render the row
		s += fmt.Sprintf("%s ⟦ %-*s ⟧: %s\n", cursor, maxFieldName, choice, value)
	}

	// The footer
	s += "\nPress s to save and quit.\nPress q to quit without saving.\n"

	// Send the UI for rendering
	return s
}
