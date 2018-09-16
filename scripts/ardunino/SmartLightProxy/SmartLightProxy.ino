#include <SerialCommand.h>
#include <SoftwareSerial.h>
#include <avr/wdt.h>

#define REDPIN   3
#define GREENPIN 6
#define BLUEPIN  9
#define WHITEPIN  10

const unsigned int MAX_INPUT = 10;
SerialCommand sCmd; 

void setup() {
  pinMode(REDPIN,   OUTPUT);
  pinMode(GREENPIN, OUTPUT);
  pinMode(BLUEPIN,  OUTPUT);
  pinMode(WHITEPIN,  OUTPUT);

  // Disable millis() and watchdog.
  TIMSK0 &= ~_BV(TOIE0);
  wdt_disable();

  // Set PWM to highest feqs possible.
  TCCR0B = TCCR0B & B11111000 | B00000001;
  TCCR1B = TCCR1B & B11111000 | B00000001;
  TCCR2B = TCCR2B & B11111000 | B00000001;

  Serial.begin(115200);
  sCmd.addCommand("SetState", SetLEDState);
  sCmd.addCommand("PWM", SetLEDPWM);
  sCmd.addCommand("Digital", SetLEDDigital);
}

void SetLEDDigital() {
  bool red, green, blue, white;
  char *arg;

  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid red boolean provided.");
    return;
  }
  red = atoi(arg) == 1;

  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid green boolean provided.");
    return;
  }
  green = atoi(arg) == 1;

  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid blue boolean provided.");
    return;
  }
  blue = atoi(arg) == 1;

  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid white boolean provided.");
    return;
  }
  white = atoi(arg) == 1;
  showDigitalRGBW(red, green, blue, white);
  
}

void SetLEDPWM() {
  uint8_t red, green, blue, white;
  char *arg;
  
  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid red colour provided.");
    return;
  }
  red = atoi(arg);

  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid green colour provided.");
    return;
  }
  green = atoi(arg);

  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid blue colour provided.");
    return;
  }
  blue = atoi(arg);

  arg = sCmd.next();
  if (arg == NULL) {
    Serial.println("ERR: Invalid white colour provided.");
    return;
  }
  white = atoi(arg);

  showAnalogRGBW(red, green, blue, white);
  
}

void SetLEDState() {
  char *arg;
  arg = sCmd.next();
  if (arg != NULL) {
    if (String(arg).equals("ON")) {
      showDigitalRGBW(true, true, true, true);
    } else if (String(arg).equals("OFF")){
      showDigitalRGBW(false, false, false, false);
    } else {
      Serial.println("ERR: invalid argument provided! Not 'ON' or 'OFF'!");
    }
    return;
  }

  Serial.println("ERR: no argument provided! Valid arguments 'ON' or 'OFF'!");
}

// Display the RGBW via PWM
void showAnalogRGBW(uint8_t red, uint8_t green, uint8_t blue, uint8_t white)
{
  analogWrite(REDPIN, red);
  analogWrite(GREENPIN, green);
  analogWrite(BLUEPIN, blue);
  analogWrite(WHITEPIN, white);
}


// Display the RGBW as solid.
void showDigitalRGBW(bool red, bool green, bool blue, bool white)
{
  if (red) {
    digitalWrite(REDPIN, HIGH );
  } else {
    digitalWrite(REDPIN, LOW );
  }

  if (green) {
    digitalWrite(REDPIN, HIGH );
  } else {
    digitalWrite(REDPIN, LOW );
  }

  if (blue) {
    digitalWrite(BLUEPIN, HIGH );
  } else {
    digitalWrite(BLUEPIN, LOW );
  }

  if (white) {
    digitalWrite(WHITEPIN, HIGH );
  } else {
    digitalWrite(WHITEPIN, LOW );
  }
}

void loop()
{
  sCmd.readSerial();
}
