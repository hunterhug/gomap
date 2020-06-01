/*
	All right reserved：https://github.com/hunterhug/gomap at 2020
	Attribution-NonCommercial-NoDerivatives 4.0 International
	You can use it for education only but can't make profits for any companies and individuals!
*/
package gomap

import (
	"fmt"
	"sync"
)

// Better AVL Tree
type avlBetterTree struct {
	c          comparator         // tree key compare
	root       *avlBetterTreeNode // tree root
	len        int64              // tree key pairs num
	sync.Mutex                    // lock for concurrent safe
}

type avlBetterTreeNode struct {
	k             string      // key
	v             interface{} // value
	left          *avlBetterTreeNode
	right         *avlBetterTreeNode
	balanceFactor int64 // balance Factor
	parent        *avlBetterTreeNode
}

// cal  height
func (node *avlBetterTreeNode) height() int64 {
	if node == nil {
		return 0
	}

	lh := node.left.height()
	rh := node.right.height()
	if lh > rh {
		return lh + 1
	} else {
		return rh + 1
	}
}

func (tree *avlBetterTree) Height() int64 {
	return tree.root.height()
}

func (tree *avlBetterTree) rotateLeft(h *avlBetterTreeNode) *avlBetterTreeNode {
	if h != nil {

		// 看图理解
		x := h.right
		h.right = x.left

		if x.left != nil {
			x.left.parent = h
		}

		x.parent = h.parent
		if h.parent == nil {
			tree.root = x
		} else if h.parent.left == h {
			h.parent.left = x
		} else {
			h.parent.right = x
		}
		x.left = h
		h.parent = x

		// can see graph
		h.balanceFactor += 1
		if x.balanceFactor < 0 {
			h.balanceFactor -= x.balanceFactor
		}

		x.balanceFactor += 1
		if h.balanceFactor > 0 {
			x.balanceFactor += h.balanceFactor
		}

		return x
	}

	return nil
}

// 对某节点右旋转
func (tree *avlBetterTree) rotateRight(h *avlBetterTreeNode) *avlBetterTreeNode {
	if h != nil {

		// 看图理解
		x := h.left
		h.left = x.right

		if x.right != nil {
			x.right.parent = h
		}

		x.parent = h.parent
		if h.parent == nil {
			tree.root = x
		} else if h.parent.right == h {
			h.parent.right = x
		} else {
			h.parent.left = x
		}
		x.right = h
		h.parent = x

		// can see graph
		h.balanceFactor -= 1
		if x.balanceFactor > 0 {
			h.balanceFactor -= x.balanceFactor
		}

		x.balanceFactor -= 1
		if h.balanceFactor < 0 {
			x.balanceFactor += h.balanceFactor
		}

		return x
	}

	return nil
}

// put key pairs
func (tree *avlBetterTree) Put(key string, value interface{}) {
	// add lock
	tree.Lock()
	defer tree.Unlock()

	if tree.root == nil {
		// 根节点都是黑色
		tree.root = &avlBetterTreeNode{
			k: key,
			v: value,
		}
		tree.len = 1
		return
	}

	var parent *avlBetterTreeNode
	var cmp int64
	node := tree.root
	for node != nil {
		cmp = tree.c(key, node.k)
		parent = node
		if cmp == 0 {
			node.v = value
			return
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	newNode := &avlBetterTreeNode{
		k:      key,
		v:      value,
		parent: parent,
	}

	if cmp < 0 {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	for parent != nil {
		// balance factor change of parent
		cmp = tree.c(parent.k, key)
		if cmp < 0 {
			parent.balanceFactor -= 1
		} else {
			parent.balanceFactor += 1
		}

		// parent factor can out, no need to loop
		if parent.balanceFactor == 0 {
			break
		} else if parent.balanceFactor < -1 {
			// right higher
			if parent.right.balanceFactor == 1 {
				tree.rotateRight(parent.right)
			}

			// root change inside
			tree.rotateLeft(parent)
			break
		} else if parent.balanceFactor > 1 {
			if parent.left.balanceFactor == -1 {
				tree.rotateLeft(parent.left)
			}
			tree.rotateRight(parent)
			break
		}

		parent = parent.parent
	}

	tree.len++
}

func (tree *avlBetterTree) Delete(key string) {
	tree.Lock()
	defer tree.Unlock()

	if tree.len == 0 {
		return
	}

	var cmp int64
	node := tree.root
	for node != nil {
		cmp = tree.c(key, node.k)
		if cmp == 0 {
			break
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	if node == nil {
		return
	}

	var maxNode, minNode *avlBetterTreeNode
	if node.left != nil {
		// find left tree max k
		maxNode = node.left

		for maxNode.left != nil || maxNode.right != nil {
			for maxNode.right != nil {
				maxNode = maxNode.right
			}

			node.k = maxNode.k
			node.v = maxNode.v

			if maxNode.left != nil {
				// max k has left node
				node = maxNode

				// let left node replace
				maxNode = maxNode.left
			}
		}

		node.k = maxNode.k
		node.v = maxNode.v
		node = maxNode // delete this node
	}

	if node.right != nil {
		minNode = node.right

		for minNode.left != nil || minNode.right != nil {
			for minNode.left != nil {
				minNode = minNode.left
			}

			node.k = minNode.k
			node.v = minNode.v

			if minNode.right != nil {
				node = minNode
				minNode = minNode.right
			}
		}

		node.k = minNode.k
		node.v = minNode.v
		node = minNode
	}

	parent := node.parent
	ps := node

	for parent != nil {
		if parent.left == ps {
			parent.balanceFactor -= 1
		} else {
			parent.balanceFactor += 1
		}

		if parent.balanceFactor < -1 {
			if parent.right.balanceFactor == 1 {
				tree.rotateRight(parent.right)
			}
			parent = tree.rotateLeft(parent)
		} else if parent.balanceFactor > 1 {
			if parent.left.balanceFactor == -1 {
				tree.rotateLeft(parent.left)
			}
			parent = tree.rotateRight(parent)
		}

		// if bal break
		if parent.balanceFactor == -1 || parent.balanceFactor == 1 {
			break
		}

		// may be continue to deal father
		ps = parent
		parent = parent.parent
	}

	if node.parent != nil {
		if node.parent.left == node {
			node.parent.left = nil
		} else {
			node.parent.right = nil
		}

		node.parent = nil
	}

	if node == tree.root {
		tree.root = nil
	}

	tree.len--
}

// find min key pairs
func (tree *avlBetterTree) MinKey() (key string, value interface{}, exist bool) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	if tree.root == nil {
		// 如果是空树，返回空
		return
	}

	node := tree.root.minNode()
	return node.k, node.v, true
}

func (node *avlBetterTreeNode) minNode() *avlBetterTreeNode {
	if node.left == nil {
		return node
	}

	return node.left.minNode()
}

// find max key pairs
func (tree *avlBetterTree) MaxKey() (key string, value interface{}, exist bool) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	if tree.root == nil {
		// 如果是空树，返回空
		return
	}

	node := tree.root.maxNode()
	return node.k, node.v, true
}

func (node *avlBetterTreeNode) maxNode() *avlBetterTreeNode {
	if node.right == nil {
		return node
	}

	return node.right.maxNode()
}

func (tree *avlBetterTree) Get(key string) (value interface{}, exist bool) {
	tree.Lock()
	defer tree.Unlock()

	if tree.root == nil {
		return
	}

	var cmp int64
	node := tree.root
	for {
		cmp = tree.c(key, node.k)
		if cmp == 0 {
			return node.v, true
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}

		if node == nil {
			break
		}
	}

	return
}

func (tree *avlBetterTree) Contains(key string) (exist bool) {
	tree.Lock()
	defer tree.Unlock()

	if tree.root == nil {
		return
	}

	var cmp int64
	node := tree.root
	for {
		cmp = tree.c(key, node.k)
		if cmp == 0 {
			return true
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}

		if node == nil {
			break
		}
	}

	return
}

func (tree *avlBetterTree) Len() int64 {
	return tree.len
}

// get int
func (tree *avlBetterTree) GetInt(key string) (value int, exist bool, err error) {
	var v interface{}
	v, exist = tree.Get(key)
	if !exist {
		return
	}

	value, ok := v.(int)
	if !ok {
		err = ReflectError(v)
		return
	}

	return value, true, nil
}

// get int64
func (tree *avlBetterTree) GetInt64(key string) (value int64, exist bool, err error) {
	var v interface{}
	v, exist = tree.Get(key)
	if !exist {
		return
	}

	value, ok := v.(int64)
	if !ok {
		err = ReflectError(v)
		return
	}

	return value, true, nil
}

// get string
func (tree *avlBetterTree) GetString(key string) (value string, exist bool, err error) {
	var v interface{}
	v, exist = tree.Get(key)
	if !exist {
		return
	}

	value, ok := v.(string)
	if !ok {
		err = ReflectError(v)
		return
	}

	return value, true, nil
}

// get float64
func (tree *avlBetterTree) GetFloat64(key string) (value float64, exist bool, err error) {
	var v interface{}
	v, exist = tree.Get(key)
	if !exist {
		return
	}

	value, ok := v.(float64)
	if !ok {
		err = ReflectError(v)
		return
	}

	return value, true, nil
}

// get byte
func (tree *avlBetterTree) GetBytes(key string) (value []byte, exist bool, err error) {
	var v interface{}
	v, exist = tree.Get(key)
	if !exist {
		return
	}

	value, ok := v.([]byte)
	if !ok {
		err = ReflectError(v)
		return
	}

	return value, true, nil
}

func (tree *avlBetterTree) KeySortedList() []string {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	keyList := make([]string, 0, tree.len)
	return tree.root.midOrder(keyList)
}

func (node *avlBetterTreeNode) midOrder(keyList []string) []string {
	if node == nil {
		return keyList
	}

	keyList = node.left.midOrder(keyList)
	keyList = append(keyList, node.k)
	keyList = node.right.midOrder(keyList)

	return keyList
}

func (tree *avlBetterTree) Check() bool {
	if tree == nil || tree.root == nil {
		return true
	}

	if tree.root.isAVL(tree.c) {
		return true
	}

	return false
}

// 判断节点是否符合 AVL 树的定义
func (node *avlBetterTreeNode) isAVL(compare comparator) bool {
	if node == nil {
		return true
	}

	if node.balanceFactor != node.left.height()-node.right.height() {
		fmt.Printf("balanceFactor %d != %d-%d\n", node.balanceFactor, node.left.height(), node.right.height())
		fmt.Printf("balanceFactor %#v != \n%#v-\n%#v\n", node, node.left, node.right)
		return false
	}
	// 左右子树都为空，那么是叶子节点
	if node.left == nil && node.right == nil {
		// 叶子节点高度应该为1
		if node.height() == 1 {
			return true
		} else {
			fmt.Println("leaf node height must is 1, now is:", node.height())
			return false
		}
	} else if node.left != nil && node.right != nil {
		// 左右子树都是满的
		// 左儿子必须比父亲小，右儿子必须比父亲大
		cmp1 := compare(node.left.k, node.k)
		cmp2 := compare(node.right.k, node.k)
		if cmp1 < 0 && cmp2 > 0 {
		} else {
			// 不符合 AVL 树定义
			fmt.Printf("father is %v lchild is %v, rchild is %v\n", node.k, node.left.k, node.right.k)
			return false
		}

		bal := node.left.height() - node.right.height()
		if bal < 0 {
			bal = -bal
		}

		// 子树高度差不能大于1
		if bal > 1 {
			fmt.Println("sub tree height bal is ", bal)
			return false
		}

		// 如果左子树比右子树高，那么父亲的高度等于左子树+1
		if node.left.height() > node.right.height() {
			if node.height() == node.left.height()+1 {
			} else {
				fmt.Printf("%#v height:%v,left sub tree height: %v,right sub tree height:%v\n", node, node.height(), node.left.height(), node.right.height())
				return false
			}
		} else {
			// 如果右子树比左子树高，那么父亲的高度等于右子树+1
			if node.height() == node.right.height()+1 {
			} else {
				fmt.Printf("%#v height:%v,left sub tree height: %v,right sub tree height:%v\n", node, node.height(), node.left.height(), node.right.height())
				return false
			}
		}

		// 递归判断子树
		if !node.left.isAVL(compare) {
			return false
		}

		// 递归判断子树
		if !node.right.isAVL(compare) {
			return false
		}

	} else {
		// 只存在一棵子树
		if node.right != nil {
			// 子树高度只能是1
			if node.right.height() == 1 && node.right.left == nil && node.right.right == nil {
				cmp := compare(node.right.k, node.k)
				if cmp > 0 {
					// 右节点必须比父亲大
				} else {
					fmt.Printf("%v has only right tree,but right small:%v", node.k, node.right.k)
					return false
				}
			} else {
				fmt.Printf("%v has only right tree height:%d,but right has lc: %#v, rc：%#v", node.k, node.right.height(), node.right.left, node.right.right)
				return false
			}
		} else {
			if node.left.height() == 1 && node.left.left == nil && node.left.right == nil {
				cmp := compare(node.left.k, node.k)
				if cmp < 0 {
					// 左节点必须比父亲小
				} else {
					fmt.Printf("%v has only left tree,but left small:%v", node.k, node.left.k)
					return false
				}
			} else {
				fmt.Printf("%v has only left tree height:%d,but left has lc: %#v, rc：%#v", node.k, node.left.height(), node.left.left, node.left.right)
				return false
			}
		}
	}

	return true
}

func (node *avlBetterTreeNode) leftOf() bsTreeNode {
	if node.left == nil {
		return nil
	}

	return node.left
}

func (node *avlBetterTreeNode) rightOf() bsTreeNode {
	if node.right == nil {
		return nil
	}

	return node.right
}

// not check node nil, may be panic, user should deal by oneself
func (node *avlBetterTreeNode) values() (key string, value interface{}) {
	return node.k, node.v
}

func (tree *avlBetterTree) KeyList() []string {
	// add lock
	tree.Lock()
	defer tree.Unlock()

	if tree.root == nil {
		return []string{}
	}

	keyList := make([]string, 0, tree.len)
	iterator := tree.Iterator()
	for iterator.HasNext() {
		k, _ := iterator.Next()
		keyList = append(keyList, k)
	}

	return keyList

}

func (tree *avlBetterTree) Iterator() MapIterator {
	q := new(linkQueue)
	if tree.root != nil {
		q.add(tree.root)
	}
	return q
}

func (tree *avlBetterTree) SetComparator(c comparator) Map {
	tree.Lock()
	defer tree.Unlock()
	if tree.len == 0 {
		tree.c = c
	}

	return tree
}

func (tree *avlBetterTree) SetHash() Map {
	tree.Lock()
	defer tree.Unlock()
	if tree.len == 0 {
		tree.c = comparatorOfSetHash
	}
	return tree
}
