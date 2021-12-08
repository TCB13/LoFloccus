#include "platformdarwin.h"
#import <Foundation/NSString.h>
#import <AppKit/AppKit.h>

void PlatformDarwin::makeAppAccessory()
{
    NSApp.activationPolicy = NSApplicationActivationPolicyAccessory;
}

void PlatformDarwin::makeAppRegular()
{
    NSApp.activationPolicy = NSApplicationActivationPolicyRegular;
}
