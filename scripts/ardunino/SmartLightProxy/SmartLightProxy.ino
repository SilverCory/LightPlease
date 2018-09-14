#include <SoftwareSerial.h>

#define REDPIN   3
#define GREENPIN 6
#define BLUEPIN  9
#define WHITEPIN  10

typedef struct {
  char instruction; // The instruction that arrived by serial connection.
  uint8_t r;
  uint8_t g;
  uint8_t b;
  uint8_t w;
} Command;

void setup() {
  pinMode(REDPIN,   OUTPUT);
  pinMode(GREENPIN, OUTPUT);
  pinMode(BLUEPIN,  OUTPUT);
  pinMode(WHITEPIN,  OUTPUT);

  // Disable millis()
  TIMSK0 &= ~_BV(TOIE0);

  // Set PWM to highest feqs possible.
  TCCR0B = TCCR0B & B11111000 | B00000001;
  TCCR1B = TCCR1B & B11111000 | B00000001;
  TCCR2B = TCCR2B & B11111000 | B00000001;

  Serial.begin(115200);
}


// Display the RGBW via PWM
void showAnalogRGBW( const Command& cmd)
{  
  analogWrite(REDPIN,   cmd.r );
  analogWrite(GREENPIN, cmd.g );
  analogWrite(BLUEPIN,  cmd.b );
  analogWrite(WHITEPIN,  cmd.w );
}


// Display the RGBW as solid.
void showDigitalRGBW( const Command& cmd)
{
  if (cmd.r == 0) {
    digitalWrite(REDPIN, LOW );
  } else {
    digitalWrite(REDPIN, HIGH );
  }

  if (cmd.g == 0) {
    digitalWrite(REDPIN, LOW );
  } else {
    digitalWrite(REDPIN, HIGH );
  }

  if (cmd.b == 0) {
    digitalWrite(BLUEPIN, LOW );
  } else {
    digitalWrite(BLUEPIN, HIGH );
  }
  
  if (cmd.w == 0) {
    digitalWrite(WHITEPIN, LOW );
  } else {
    digitalWrite(WHITEPIN, HIGH );
  }
}

void loop()
{
  Command cmd = ReadCommand();
  if (cmd.instruction == 'P') {
    showAnalogRGBW(cmd);
  }
//  } else if (cmd.instruction == 'D') {
//    showDigitalRGBW(cmd);
//  }
}

/**
 * ReadCommand sucks down the lastest command from the serial port,
 * returns {'*', 0.0} if no new command is available.
 */
Command ReadCommand() {
  // Not enough bytes for a command, return an empty command.
  if (Serial.available() < 5) {
    return (Command) {'*', 0.0};
  }
  
  char c = Serial.read();
  uint8_t r = Serial.read();
  uint8_t g = Serial.read();
  uint8_t b = Serial.read();
  uint8_t w = Serial.read();

  return (Command) {c, r, g, b, w};
}


