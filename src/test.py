import copy

class Foo:
    bar = None
    val = 0

class Bar:
    foo = None
    val = 0

f1 = Foo()
b1 = Bar()
f1.bar = b1
b1.foo = f1
b1.val = 1
f1.val = 1

cp = copy.deepcopy(f1)
b1.val = 2
f1.val = 2

print("Copy foo = " + str(cp.val))
print("Copy bar = " + str(cp.bar.val))
print("Copy foo from bar = " + str(cp.bar.foo.val))