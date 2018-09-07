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

void loop()
{
  Command cmd = ReadCommand();
  if (cmd.instruction == 'F') {
    showAnalogRGBW(cmd);
  }
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


