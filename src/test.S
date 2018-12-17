.text

main:
jal read
nop
li $a1,str_test
li $a2,0x0800
jal streq
nop
beq $v0,$zero,main
nop
main_end:
addi $v0,$zero,10
syscall

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
li $v0,10
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
li $a0,0x0800
li $a1,0x0010
li $v0,8
syscall
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