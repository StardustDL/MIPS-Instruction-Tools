.text

main:
jal read
lui $a1,str_test
lui $a2,str_hello
addiu $a1,$a1,0x1800
addiu $a2,$a2,0
jal streq
nop
bne $v0,$zero,hello_world
nop
hello_world:
lui $a1,0x00000000
addiu $a1,$a1,0x1800
jal write
nop
j main
nop

streq: # compare string s1 and s2
streq_loop:
lb $t1,($a1)
lb $t2,($a2)
or $t3,$t1,$t2
beq $t3,$zero,streq_eq
nop
bne $t1,$t2,streq_ne
nop
addiu $a1,$a1,1
addiu $a2,$a2,1
j streq_loop
nop
streq_ne: 
addiu $v0,$zero,0
jr $ra
nop
streq_eq:
addiu $v0,$zero,1
jr $ra
nop

strcpy:
strcpy_loop:
lb $t1,($a1)
sb $t1,($a2)
beq $t1,$zero,strcpy_end
nop
addiu $a1,$a1,1
addiu $a2,$a2,1
j strcpy_loop
nop
strcpy_end:
jr $ra
nop

read:
read_before:
lb $t1,0x0000
beq $t1,$zero,read_before
nop
j read_main
nop
read_main:
lui $a1,0x0000
lui $a2,0x1800
jal strcpy
nop
sb $zero,0x0000
jr $ra
nop

write:
lui $a2,0x0080
jal strcpy
nop
jr $ra
nop

.data
str_test: .asciiz "test"
str_hello: .asciiz "hello world\r"
str_abc: .asciiz "abc"