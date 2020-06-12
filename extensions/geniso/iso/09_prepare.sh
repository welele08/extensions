#!/bin/sh

set -ex

echo "**** Copy kernel"
# Prepare the kernel install area.
echo "Removing old kernel artifacts. This may take a while."
rm -rf $KERNEL_INSTALLED
mkdir -p $KERNEL_INSTALLED
BOOT_DIR=$ROOTFS_DIR/boot
if [[ -L "$BOOT_DIR/bzImage" ]]
then

bz=$(readlink -f $BOOT_DIR/bzImage)
# Install the kernel file.
cp $BOOT_DIR/$(basename $bz) \
  $KERNEL_INSTALLED/kernel
else 
cp $BOOT_DIR/bzImage \
  $KERNEL_INSTALLED/kernel
fi
echo "*** GENERATE ROOTFS BEGIN ***"


