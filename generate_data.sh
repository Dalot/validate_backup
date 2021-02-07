#!/bin/bash

array=(NoFe0qnuu6 snRDxBxbmM vqUSSLqldG FI1nIFv3Ys LmiXlP21oW Fx7RBcDRSB mk0pfy0bua Q7PI34p8wM fH3JfC56ig F2IxqWYWLU Buo56XJopP Yd2jozImH4 hjV5LzpFhN IvWIW8Ea7w G6uo3w1RKH U4HzKMLzWU tQBSEoJYd3 8uPFmq3UZV eUhsyiUjsQ l9xXUYiQkv SCWu5w8Zn7 OKaRfgtR3v HZUsPpPIdE 2yayJECee4 NbkfJrCn9K hc5mr3Ykp9 uNyLG8Si8b GsW2wAd796 edzJiJ3ZNQ T3J3xSsbDH WW1oPKhdsv eEl8QtXKXl HHdxyHivsL Eu42OjC8rP 8ol6jwMRhk 6yoLSsK9NI IgIbeSK6G7 a3wfhGAEpl sdoj1ex1tN dOzsorXVdx 9UsdqD1KId wmnv83vyun D87mbzUbpQ 1g88y5xbjE 9vO2iLxbYU eIwaxN54En 4xghT245Nr M6wXkGGgmi UVsBJqk5H1 FYTuXaIedV un2GBvd3OU 54EzMLzMbj Zg1HlNZE76 3bbbo3zAem bind9HhAh6 6yXuesrhXp 1KdGBgT4Lu tFfz2mjgCo z0O8cy1xwF F3ddSLl0yi t75UuglTjO mifTHOVGv4 6sB4jPmy1B 8wcwzglVOE 8zDor7rAGM DhKA6Ollhb cse01ss6u7 PvtXUejdJ0 Vprn48dxJk 0ngn3foM82 isrrNaFj7B a7IeZNQ2pl VBLSZzMqx5 dBTynwaW2D yHXdEpipBX TirUmwOri5 3b8DGcII1U vBr9UpVK4q DRWggIK3x2 VXBmCJLI3C NBpSIfxnRZ lSxSgncidm NIClD6jAkU NKqQortaBq CTWSZp72DI g0CpHfUHbk XDAWMIrDSG mgsfq8NUT7 oMuzD6G2Gz x7kcmzCMMX 5T8hremXJD uRekPxDsOx u2Km4CGVET eLhGmG7Agv qaHmE91X3N YruuhIDDvg nBL2zleLHI hViHKOzQcu CPLdxLdDxZ CWQEMNwY8q)

file=./xl_after.json
#echo [ > $file
start=`date +%s`
for i in `seq 7000001 10000000`;
    do
        let "id1=$i+1"
        let "id2=$i+2"
        let "id3=$i+3"
        let "id4=$i+4"
        let "id5=$i+5"
        let "id6=$i+6"
        let "id7=$i+7"
        let "id8=$i+8"
        let "id9=$i+9"
        let "id10=$i+10"
        let "id11=$i+11"
        let "id12=$i+12"
        let "id13=$i+13"
        let "id14=$i+14"
        let "id15=$i+15"
        let "id16=$i+16"
        let "id17=$i+17"
        let "id18=$i+18"
        let "id19=$i+19"
        let "id20=$i+20"
        #size=${#array[@]}
        #index=$(($RANDOM % $size))
        #$name=${array[$index]}
        name="hatch_challenge"
        name1=$name$id1
        name2=$name$id2
        name3=$name$id3
        name4=$name$id4
        name5=$name$id5
        name6=$name$id6
        name7=$name$id7
        name8=$name$id8
        name9=$name$id9
        name10=$name$id10
        name11=$name$id11
        name12=$name$id12
        name13=$name$id13
        name14=$name$id14
        name15=$name$id15
        name16=$name$id16
        name17=$name$id17
        name18=$name$id18
        name19=$name$id19
        name20=$name$id20
        
        elems="{\"id\": \"$i\", \"name\": \"$name\"},{\"id\": \"$id1\", \"name\": \"$name1\"},{\"id\": \"$id2\", \"name\": \"$name2\"},{\"id\": \"$id3\", \"name\": \"$name3\"},{\"id\": \"$id4\", \"name\": \"$name4\"},{\"id\": \"$id5\", \"name\": \"$name5\"},{\"id\": \"$id6\", \"name\": \"$name6\"},{\"id\": \"$id7\", \"name\": \"$name7\"},{\"id\": \"$id8\", \"name\": \"$name8\"},{\"id\": \"$id9\", \"name\": \"$name9\"},{\"id\": \"$id10\", \"name\": \"$name10\"},{\"id\": \"$id11\", \"name\": \"$name11\"},{\"id\": \"$id12\", \"name\": \"$name12\"},{\"id\": \"$id13\", \"name\": \"$name13\"},{\"id\": \"$id14\", \"name\": \"$name14\"},{\"id\": \"$id15\", \"name\": \"$name15\"},{\"id\": \"$id16\", \"name\": \"$name16\"},{\"id\": \"$id17\", \"name\": \"$name17\"},{\"id\": \"$id18\", \"name\": \"$name18\"},{\"id\": \"$id19\", \"name\": \"$name19\"},{\"id\": \"$id20\", \"name\": \"$name20\"},"
        elems="{\"id\": \"$id1\", \"name\": \"$name1\"},"
        echo $elems >> $file
        i=$i+20
    done
echo "{\"id\": \"final\", \"name\": \"countdown\"}" >> $file
echo ] >> $file
end=`date +%s`
echo "time taken: $((end-start))"
