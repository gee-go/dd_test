package ddlog

// Parser is used to convert lines to Message structs
type Parser struct {
	b      *buffer
	config *Config
}

func NewParser(c *Config) *Parser {
	return &Parser{config: c, b: &buffer{}}
}

// Parse converts a line to a message.
func (p *Parser) Parse(l string) (*Message, error) {
	p.b.Init(l)

	end := ' '
	fieldStart := -1
	prev := ' '
	msg := &Message{}

	for i, r := range p.config.LogFormat {
		// field names
		switch r {
		case '{':
			fieldStart = i + 1
			switch prev {
			case '[':
				end = ']'
			case '"':
				end = '"'
			default:
				end = prev
			}
		case '}':
			if err := msg.set(p.config.LogFormat[fieldStart:i], p.b.advanceUntil(end)); err != nil {
				return msg, err
			}
		}

		prev = r
	}

	return msg, nil
}
