#!/bin/bash

set -e

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/func.sh"

XZ="${XZ:-xz}"


if [[ -z "$INITRAMFS_ROOTFS" ]]; then
    info "Packing initramfs"

    # Remove the old 'initramfs' archive if it exists.
    rm -f $WORKDIR/rootfs.cpio.xz || true

    pushd $ROOTFS_DIR  > /dev/null 2>&1

    # Packs the current 'initramfs' folder structure in 'cpio.xz' archive.
    find . | cpio -R root:root -H newc -o | $XZ -9 --check=none > $WORKDIR/rootfs.cpio.xz

    echo "Packing of initramfs has finished."

    popd  > /dev/null 2>&1
else
    echo "Copying initramfs"
    # Try to find the rootfs file in the overlay or initramfs areas
    if [[ -e "$ROOTFS_DIR/boot/$INITRAMFS_ROOTFS" ]] || [[ -L "$ROOTFS_DIR/boot/$INITRAMFS_ROOTFS" ]]; then
        BOOT_DIR=$ROOTFS_DIR/boot
    elif [[ -e "$OVERLAY_DIR/boot/$INITRAMFS_ROOTFS" ]] || [[ -L "$OVERLAY_DIR/boot/$INITRAMFS_ROOTFS" ]]; then
        BOOT_DIR=$OVERLAY_DIR/boot
    fi

    if [[ -L "$BOOT_DIR/$INITRAMFS_ROOTFS" ]]; then
        bz=$(readlink -f $BOOT_DIR/$INITRAMFS_ROOTFS)
        # Install the kernel file.
        cp $BOOT_DIR/$(basename $bz) \
            $WORKDIR/rootfs.cpio.xz
    else
        cp $BOOT_DIR/$INITRAMFS_ROOTFS \
            $WORKDIR/rootfs.cpio.xz 
    fi
fi

info "Packing overlayfs"
rm -f $WORKDIR/rootfs.squashfs || true
mksquashfs "$OVERLAY_DIR" $WORKDIR/rootfs.squashfs -b 1024k -comp xz -Xbcj x86
