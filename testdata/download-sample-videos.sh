#!/bin/sh

DIR_ROOT1=./root1
DIR_ROOT2=./root2
DIR_ROOT2_SUB1=$DIR_ROOT2/root2-sub1

MP4_URL=https://www.sample-videos.com/video321/mp4/240/big_buck_bunny_240p_1mb.mp4
MKV_URL=https://www.sample-videos.com/video321/mkv/240/big_buck_bunny_240p_1mb.mkv
AVI_URL=https://file-examples.com/storage/fe9a194958686838db9645f/2018/04/file_example_AVI_480_750kB.avi

mkdir -p $DIR_ROOT1
mkdir -p $DIR_ROOT2
mkdir -p $DIR_ROOT2_SUB1

curl -o "$DIR_ROOT1/root1-01.mp4" $MP4_URL
curl -o "$DIR_ROOT1/root1-02.avi" $AVI_URL
curl -o "$DIR_ROOT2/root2-01.mp4" $MP4_URL
curl -o "$DIR_ROOT2/root2-02.mkv" $MKV_URL
curl -o "$DIR_ROOT2/root2-03.mp4" $MP4_URL
curl -o "$DIR_ROOT2_SUB1/root2-sub1-01.mp4" $MP4_URL
