#include "anykey/anykey.h"
#include "anykey_usb/keyboard.h"

#define LED_PORT 0
#define LED_PIN 7
#define KEY_PORT 1
#define KEY_PIN 1

#define RGB_PORT 0
#define R_PIN 8
#define G_PIN 9
#define B_PIN 10

#define KEY_DOWN_DURATION 25 /* time in 10ms-ticks */
#define LIGHT_ON_DURATION 500 /* time in 10ms-ticks */

/** these variables are needed by the Keyboard implementation. USB peripherals don't allocate
 memory by themselves, so we have to give some memory to them. They usually handle all related stuff. */
USB_Device_Definition usbDeviceDefinition;
USB_Device_Struct usbDevice;
USBHID_Behaviour_Struct hidBehaviour;
uint8_t inBuffer[8];
uint8_t outBuffer[8];
uint8_t idleValue;
uint8_t currentProtocol;

/** remember when the button was down */
bool buttonWasDown = false; 

typedef enum {
	State_AllOff = 0,
	State_P1_Enabled = 1,
	State_P2_Enabled = 2,
	State_P1_P2_Enabled = 3,
	State_P3_On_Enabled = 4,
	State_P1_P3_Enabled = 5,
	State_P2_P3_Enabled = 6,
	State_P1_P2_P3_Enabled = 7,
	State_P1_Only = 8,
	State_P2_Only = 9,
	State_P3_Only = 10,
	State_P1_Urgent = 11,
	State_P2_Urgent = 12,
	State_P3_Urgent = 13,
	State_Rainbow = 14,
	State_Alone = 32				//Key was pressed, own color is shown, will go back to allOff after some delay. Used when no feedback is given
} ButtonState;

ButtonState state = State_AllOff;

/** indicates "our" button ID (red, green or blue) */
const uint8_t myButtonIdx = 1;

/** 10-ms ticks in current state */
uint32_t stateTicks = 0;

/** 10-ms for typing. Negative = not typing */
int typeTicks = -1;

/** total number of buttons */
const uint8_t numButtons = 3;

/** USB HID keyboard codes. See HUT1_12.pdf (from usb.org). */
const uint8_t keys[] = { 21, 10, 5};

typedef struct {
	uint16_t r;
	uint16_t g;
	uint16_t b;
} RGBColor;

/** colors for each button */
const RGBColor colors[] = {
	{0xffff, 0, 0},
	{0, 0xffff, 0},
	{0, 0, 0xffff},
};

const RGBColor rainbowColors[] = {
	{0xffff, 0, 0},
	{0, 0xffff, 0},
	{0, 0, 0xffff},
	{0xffff, 0, 0},
	{0, 0xffff, 0},
	{0, 0, 0xffff},

	{0xffff, 0, 0},
	{0xffff, 0xffff, 0},
	{0, 0xffff, 0},
	{0, 0xffff, 0xffff},
	{0, 0, 0xffff},
	{0xffff, 0, 0xffff},
};


const RGBColor black = {0,0,0};

/** USB HID keyboard codes. See HUT1_12.pdf (from usb.org). */
uint8_t downKey = 0;

/** generate an IN report. We only use one key, no modifiers */
uint16_t inReportHandler (USB_Device_Struct* device,
						  const USBHID_Behaviour_Struct* behaviour,
						  USB_HID_REPORTTYPE reportType,
						  uint8_t reportId) {
	inBuffer[0] = 0;
	inBuffer[1] = 0;
	inBuffer[2] = downKey;
	inBuffer[3] = 0;
	inBuffer[4] = 0;
	inBuffer[5] = 0;
	inBuffer[6] = 0;
	inBuffer[7] = 0;
	return 8;
}

/** parse an OUT report. We just read the caps lock bit and turn the LED on and
 * off */

void outReportHandler (USB_Device_Struct* device,
					   const USBHID_Behaviour_Struct* behaviour,
					   USB_HID_REPORTTYPE reportType,
					   uint8_t reportId,
					   uint16_t len) {
	state = outBuffer[0];
	stateTicks = 0;
}

void SendKey(uint8_t key) {
	downKey = key;
	USBHID_PushReport (&usbDevice, &hidBehaviour, USB_HID_REPORTTYPE_INPUT, 0);
}

void SetRGBColor(const RGBColor* color) {
	Timer_SetMatchValue(CT16B0, 0, ~color->r);
	Timer_SetMatchValue(CT16B0, 1, ~color->g);
	Timer_SetMatchValue(CT16B0, 2, ~color->b);
}

void systick () {
	stateTicks++;
	if (typeTicks >= 0) typeTicks++;

	//Check if button was pressed. TODO: Disable button events in certain states
	bool buttonDown = !any_gpio_read (KEY_PORT, KEY_PIN);
	if (buttonDown && !buttonWasDown && (typeTicks < 0)) {
		typeTicks = 0;
		if ((state == State_AllOff) || (state == State_Alone)) {
			state = State_Alone;
			stateTicks = 0;
		} 
	}
	buttonWasDown = buttonDown;

	//Special case no feedback: If we're all alone, go back to alloff after a while
	if ((state == State_Alone) && (stateTicks > LIGHT_ON_DURATION)) {
		state = State_AllOff;
		stateTicks = 0;
	}
	
	//Send typing HID events.
	if (typeTicks == 0) {
		SendKey(keys[myButtonIdx]);
	} else if (typeTicks == KEY_DOWN_DURATION) {
		SendKey(0);
		typeTicks = -1;
	}

	//Set RGB color according to state
	if (state <= State_P1_P2_P3_Enabled) {
		bool on = state & (1 << myButtonIdx);
		SetRGBColor(on ? &colors[myButtonIdx] : &black);
	} else if ((state >= State_P1_Only) && (state <= State_P3_Urgent)) {
		bool blink = state >= State_P1_Urgent;
		int activeButtonIdx = (state - State_P1_Only) % numButtons;
		bool on = blink ? ((stateTicks & 0x10) == 0) : true;
		SetRGBColor(on ? &colors[activeButtonIdx] : &black);
	} else if (state == State_Rainbow) {
		int fromIdx = (stateTicks >> 7) % 6;
		int toIdx = (fromIdx + 1) % 6;
		uint32_t weight_to   = stateTicks & 0x7f;
        uint32_t weight_from = 0x80 - weight_to;
        uint32_t r = (weight_from * rainbowColors[fromIdx].r + weight_to * rainbowColors[toIdx].r ) >> 7;
        uint32_t g = (weight_from * rainbowColors[fromIdx].g + weight_to * rainbowColors[toIdx].g ) >> 7;
        uint32_t b = (weight_from * rainbowColors[fromIdx].b + weight_to * rainbowColors[toIdx].b ) >> 7;
        RGBColor blendColor = { r,g,b };
		SetRGBColor(&blendColor);
	} else if (state == State_Alone) {
		SetRGBColor(&colors[myButtonIdx]);
	}
}

void main () {
	any_gpio_set_dir (LED_PORT, LED_PIN, OUTPUT);
	any_gpio_write (LED_PORT, LED_PIN, false);
	any_gpio_set_dir (KEY_PORT, KEY_PIN, INPUT);
	ANY_GPIO_SET_PULL (KEY_PORT, KEY_PIN, PULL_UP);
	
	// prepare the RGB outputs for PWM (timer and digital mode)
	ANY_GPIO_SET_FUNCTION(RGB_PORT, R_PIN, TMR, IOCON_IO_ADMODE_DIGITAL);
	ANY_GPIO_SET_FUNCTION(RGB_PORT, G_PIN, TMR, IOCON_IO_ADMODE_DIGITAL);
	ANY_GPIO_SET_FUNCTION(RGB_PORT, B_PIN, TMR, IOCON_IO_ADMODE_DIGITAL);

	// enable the timer
	Timer_Enable(CT16B0, true);

	SetRGBColor(&black);

	Timer_SetMatchBehaviour(CT16B0, 0, 0);
	// enable PWM on red pin ...
	Timer_EnablePWM(CT16B0, 0, true);

	// repeat for green ...
	Timer_SetMatchBehaviour(CT16B0, 1, 0);
	Timer_EnablePWM(CT16B0, 1, true);

	// ... and blue
	Timer_SetMatchBehaviour(CT16B0, 2, 0);
	Timer_EnablePWM(CT16B0, 2, true);

	Timer_SetMatchBehaviour(CT16B0, 3, 0);

	Timer_Start(CT16B0);

	KeyboardInit (&usbDeviceDefinition,
				  &usbDevice,
				  &hidBehaviour,
				  inBuffer,
				  outBuffer,
				  &idleValue,
				  &currentProtocol,
				  inReportHandler,
				  outReportHandler);
	
	USB_SoftConnect (&usbDevice);

	SYSCON_StartSystick_10ms();
}
