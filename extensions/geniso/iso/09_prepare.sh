#!/bin/bash

set -e

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/func.sh"

info "Copy kernel"

rm -rf $KERNEL_INSTALLED || true
mkdir -p $KERNEL_INSTALLED

# Try to find the kernel file in the overlay or initramfs areas
if [[ -e "$ROOTFS_DIR/boot/$INITRAMFS_KERNEL" ]] || [[ -L "$ROOTFS_DIR/boot/$INITRAMFS_KERNEL" ]]; then
  BOOT_DIR=$ROOTFS_DIR/boot
elif [[ -e "$OVERLAY_DIR/boot/$INITRAMFS_KERNEL" ]] || [[ -L "$OVERLAY_DIR/boot/$INITRAMFS_KERNEL" ]]; then
  BOOT_DIR=$OVERLAY_DIR/boot
fi

if [[ -L "$BOOT_DIR/$INITRAMFS_KERNEL" ]]; then
  bz=$(readlink -f $BOOT_DIR/$INITRAMFS_KERNEL)
  # Install the kernel file.
  cp $BOOT_DIR/$(basename $bz) \
    $KERNEL_INSTALLED/kernel
else
  cp $BOOT_DIR/$INITRAMFS_KERNEL \
    $KERNEL_INSTALLED/kernel
fi
