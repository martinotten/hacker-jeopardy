// LED setter

#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/hid/IOHIDLib.h>
#import <Cocoa/Cocoa.h>
uint8_t setMask;

void setLEDs(const void *value, void *context);

int main( int argc, const char * argv[] ) {
	if (argc != 2) {
        printf("Usage: setleds <bitmask as int>\n");
        return 0;
    }
    setMask = atoi(argv[1]);

    NSAutoreleasePool* pool = [[NSAutoreleasePool alloc] init];
    
	IOHIDManagerRef hidManager = IOHIDManagerCreate( kCFAllocatorDefault, kIOHIDOptionsTypeNone );
	if (!hidManager) return 1;
	
    NSDictionary* keyboardMatchDict = [NSDictionary dictionaryWithObjectsAndKeys:
                                       [NSNumber numberWithInt: kHIDPage_GenericDesktop],@ kIOHIDDeviceUsagePageKey,
                                       [NSNumber numberWithInt: kHIDUsage_GD_Keyboard],@ kIOHIDDeviceUsageKey, nil];
    
	if (!keyboardMatchDict) return 2;

	IOHIDManagerSetDeviceMatching( hidManager, (CFDictionaryRef) keyboardMatchDict);
	
	IOReturn status = IOHIDManagerOpen( hidManager, kIOHIDOptionsTypeNone );
    if (status) return 3;
	
	CFSetRef keyboards = IOHIDManagerCopyDevices(hidManager);
    if (!keyboards) return 4;

    CFSetApplyFunction(keyboards, setLEDs, NULL);
	
    CFRelease(hidManager);
    [pool drain];
    return 0;
}

void setLEDs(const void *value, void *context) {
#pragma unused(context)
    
    IOHIDDeviceRef device = (IOHIDDeviceRef)value;
    if (!IOHIDDeviceConformsTo(device, kHIDPage_GenericDesktop, kHIDUsage_GD_Keyboard)) return;
    IOHIDDeviceSetReport(device, kIOHIDReportTypeOutput, 0, &setMask, 1);
    
}
