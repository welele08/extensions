#!/bin/sh

set -ex

XZ="${XZ:-xz}"

echo "*** PACK ROOTFS BEGIN ***"

echo "Packing initramfs. This may take a while."

# Remove the old 'initramfs' archive if it exists.
rm -f $WORKDIR/rootfs.cpio.xz

pushd $ROOTFS_DIR

# Packs the current 'initramfs' folder structure in 'cpio.xz' archive.
find . | cpio -R root:root -H newc -o | $XZ -9 --check=none > $WORKDIR/rootfs.cpio.xz

echo "Packing of initramfs has finished."

popd
rm -f $WORKDIR/rootfs.squashfs
mksquashfs "$OVERLAY_DIR" $WORKDIR/rootfs.squashfs -b 1024k -comp xz -Xbcj x86