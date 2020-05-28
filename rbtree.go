/*
	All right reserved：https://github.com/hunterhug/gomap at 2020
	Attribution-NonCommercial-NoDerivatives 4.0 International
	You can use it for education only but can't make profits for any companies and individuals!
*/
package gomap

import (
	"errors"
	"fmt"
	"sync"
)

// color
const (
	RED   = true
	BLACK = false
)

// type reflect err
func ReflectError(v interface{}) error {
	return errors.New(fmt.Sprintf("type is %T, value is:%#v", v, v))
}

// red-black tree, short call rbt
// refer Java TreeMap
// the first version is comment by chinese so keep it, future new comment will be Eng
type rbTree struct {
	root       *rbTNode // tree root node
	len        int64    // tree key pairs num
	sync.Mutex          // lock for concurrent safe
}

// new rbt tree map
func NewRBTreeMap() Map {
	return new(rbTree)
}

// rbt node
// all field is lowercase to keep black box
type rbTNode struct {
	k      string      // key
	v      interface{} // value
	left   *rbTNode    // left tree
	right  *rbTNode    // right tree
	parent *rbTNode    // node's parent
	color  bool        // color of parent point to this node
}

// is rbt node is red
func isRed(node *rbTNode) bool {
	if node == nil {
		return false
	}
	return node.color == RED
}

// 返回节点的父亲节点
func parentOf(node *rbTNode) *rbTNode {
	if node == nil {
		return nil
	}

	return node.parent
}

// 返回节点的左子节点
func leftOf(node *rbTNode) *rbTNode {
	if node == nil {
		return nil
	}

	return node.left
}

// 返回节点的右子节点
func rightOf(node *rbTNode) *rbTNode {
	if node == nil {
		return nil
	}

	return node.right
}

// 设置节点颜色
func setColor(node *rbTNode, color bool) {
	if node != nil {
		node.color = color
	}
}

// 对某节点左旋转
func (tree *rbTree) rotateLeft(h *rbTNode) {
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
	}
}

// 对某节点右旋转
func (tree *rbTree) rotateRight(h *rbTNode) {
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
	}
}

// 普通红黑树添加元素
// put key pairs
func (tree *rbTree) Put(key string, value interface{}) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	//fmt.Println("add,", key)

	// 根节点为空
	if tree.root == nil {
		// 根节点都是黑色
		tree.root = &rbTNode{
			k:     key,
			v:     value,
			color: BLACK,
		}
		tree.len = 1
		return
	}

	// 辅助变量 t，表示新元素要插入到该子树，t是该子树的根节点
	t := tree.root

	// 插入元素后，插入元素的父亲节点
	var parent *rbTNode

	// 辅助变量，为了知道元素最后要插到左边还是右边
	var cmp int64 = 0

	for {
		parent = t

		cmp = compare(key, t.k)
		if cmp < 0 {
			// 比当前节点小，往左子树插入
			t = t.left
		} else if cmp > 0 {
			// 比当前节点节点大，往右子树插入
			t = t.right
		} else {
			// update new value
			t.v = value
			return
		}

		// 终于找到要插入的位置了
		if t == nil {
			break // 这时叶子节点是 parent，要插入到 parent 的下面，跳到外层去
		}
	}

	// 新节点，它要插入到 parent下面
	newNode := &rbTNode{
		k:      key,
		v:      value,
		parent: parent,
	}
	if cmp < 0 {
		// 知道要从左边插进去
		parent.left = newNode
	} else {
		// 知道要从右边插进去
		parent.right = newNode
	}

	// 插入新节点后，可能破坏了红黑树特征，需要修复，核心函数
	tree.fixAfterInsertion(newNode)

	// len add 1
	tree.len = tree.len + 1
}

// 调整新插入的节点，自底而上
// 可以看图理解
func (tree *rbTree) fixAfterInsertion(node *rbTNode) {
	// 插入的新节点一定要是红色
	node.color = RED

	// 节点不能是空，不能是根节点，父亲的颜色必须为红色（如果是黑色，那么直接插入不破坏平衡，不需要调整了）
	for node != nil && node != tree.root && node.parent.color == RED {
		// 父亲在祖父的左边
		if parentOf(node) == leftOf(parentOf(parentOf(node))) {
			// 叔叔节点
			uncle := rightOf(parentOf(parentOf(node)))

			// 图例3左边部分，叔叔是红节点，祖父变色，也就是父亲和叔叔变黑，祖父变红
			if isRed(uncle) {
				setColor(parentOf(node), BLACK)
				setColor(uncle, BLACK)
				setColor(parentOf(parentOf(node)), RED)
				// 还要向上递归
				node = parentOf(parentOf(node))
			} else {
				// 图例4左边部分，叔叔是黑节点，并且插入的节点在父亲的右边，需要对父亲左旋
				if node == rightOf(parentOf(node)) {
					node = parentOf(node)
					tree.rotateLeft(node)
				}

				// 变色，并对祖父进行右旋
				setColor(parentOf(node), BLACK)
				setColor(parentOf(parentOf(node)), RED)
				tree.rotateRight(parentOf(parentOf(node)))
			}
		} else {
			// 父亲在祖父的右边，与父亲在祖父的左边相似
			// 叔叔节点
			uncle := leftOf(parentOf(parentOf(node)))

			// 图例3右边部分，叔叔是红节点，祖父变色，也就是父亲和叔叔变黑，祖父变红
			if isRed(uncle) {
				setColor(parentOf(node), BLACK)
				setColor(uncle, BLACK)
				setColor(parentOf(parentOf(node)), RED)
				// 还要向上递归
				node = parentOf(parentOf(node))
			} else {
				// 图例4右边部分，叔叔是黑节点，并且插入的节点在父亲的左边，需要对父亲右旋
				if node == leftOf(parentOf(node)) {
					node = parentOf(node)
					tree.rotateRight(node)
				}

				// 变色，并对祖父进行左旋
				setColor(parentOf(node), BLACK)
				setColor(parentOf(parentOf(node)), RED)
				tree.rotateLeft(parentOf(parentOf(node)))
			}
		}
	}

	// 根节点永远为黑
	tree.root.color = BLACK
}

// 普通红黑树删除元素
func (tree *rbTree) Delete(key string) {
	tree.Lock()
	defer tree.Unlock()

	if tree.len == 0 {
		return
	}
	// 查找元素是否存在，不存在则退出
	node := tree.find(key)
	if node == nil {
		return
	}

	//fmt.Println("delete,", key)
	// 删除该节点
	tree.delete(node)

	tree.len = tree.len - 1
}

// 删除节点核心函数
// 找最小后驱节点来补位，删除内部节点转为删除叶子节点
func (tree *rbTree) delete(node *rbTNode) {
	// 如果左右子树都存在，那么从右子树的左边一直找一直找，就找能到最小后驱节点
	if node.left != nil && node.right != nil {
		s := node.right
		for s.left != nil {
			s = s.left
		}

		// 删除的叶子节点找到了，删除内部节点转为删除叶子节点
		node.k = s.k
		node.v = s.v
		node = s // node may be has one right son
	}

	if node.left == nil && node.right == nil {
		// 没有子树，要删除的节点就是叶子节点。
	} else {
		// 只有一棵子树，因为红黑树的特征，该子树就只有一个节点
		// 找到该唯一节点
		replacement := node.left
		if node.left == nil {
			replacement = node.right
		}

		// 替换开始，子树的唯一节点替代被删除的内部节点
		replacement.parent = node.parent

		if node.parent == nil {
			// 要删除的节点的父亲为空，表示要删除的节点为根节点，唯一子节点成为树根
			tree.root = replacement
		} else if node == node.parent.left {
			// 子树的唯一节点替代被删除的内部节点
			node.parent.left = replacement
		} else {
			// 子树的唯一节点替代被删除的内部节点
			node.parent.right = replacement
		}

		// delete this node
		node.parent = nil
		node.right = nil
		node.left = nil

		//  case 1: not enter this logic
		//      R(del)
		//    B   B
		//
		//  case 2: node's color must be black, and it's son must be red
		//    B(del)     B(del)
		//  R  O       O   R
		//
		// 单子树时删除的节点绝对是黑色的，而其唯一子节点必然是红色的
		// 现在唯一子节点替换了被删除节点，该节点要变为黑色
		// now son replace it's father, just change color to black
		replacement.color = BLACK

		//// 要删除的节点，是一个黑节点，删除后会破坏平衡，需要进行调整，调整成可以删除的状态
		//if !isRed(node) {
		//	// 核心函数
		//	tree.fixAfterDeletion(replacement)
		//}

		return
	}

	// 要删除的叶子节点没有父亲，那么它是根节点，直接置空，返回
	if node.parent == nil {
		//fmt.Println("root")
		tree.root = nil
		return
	}

	// 要删除的叶子节点，是一个黑节点，删除后会破坏平衡，需要进行调整，调整成可以删除的状态
	if !isRed(node) {
		// 核心函数
		tree.fixAfterDeletion(node)
	}

	// 现在可以删除叶子节点了
	if node == node.parent.left {
		node.parent.left = nil
	} else if node == node.parent.right {
		node.parent.right = nil
	}

	node.parent = nil

}

// 调整删除的叶子节点，自底向上
// 可以看图理解
func (tree *rbTree) fixAfterDeletion(node *rbTNode) {
	// 如果不是递归到根节点，且节点是黑节点，那么继续递归
	for tree.root != node && !isRed(node) {
		// 要删除的节点在父亲左边，对应图例1，2
		if node == leftOf(parentOf(node)) {
			// 找出兄弟
			brother := rightOf(parentOf(node))

			// 兄弟是红色的，对应图例1，那么兄弟变黑，父亲变红，然后对父亲左旋，进入图例21,22,23
			if isRed(brother) {
				setColor(brother, BLACK)
				setColor(parentOf(node), RED)
				tree.rotateLeft(parentOf(node))
				brother = rightOf(parentOf(node)) // 图例1调整后进入图例21,22,23，兄弟此时变了
			}

			// 兄弟是黑色的，对应图例21，22，23
			// 兄弟的左右儿子都是黑色，进入图例23，将兄弟设为红色，父亲所在的子树作为整体，当作删除的节点，继续向上递归
			if !isRed(leftOf(brother)) && !isRed(rightOf(brother)) {
				setColor(brother, RED)
				node = parentOf(node)
			} else {
				// 兄弟的右儿子是黑色，进入图例22，将兄弟设为红色，兄弟的左儿子设为黑色，对兄弟右旋，进入图例21
				if !isRed(rightOf(brother)) {
					setColor(leftOf(brother), BLACK)
					setColor(brother, RED)
					tree.rotateRight(brother)
					brother = rightOf(parentOf(node)) // 图例22调整后进入图例21，兄弟此时变了
				}

				// 兄弟的右儿子是红色，进入图例21，将兄弟设置为父亲的颜色，兄弟的右儿子以及父亲变黑，对父亲左旋
				setColor(brother, parentOf(node).color)
				setColor(parentOf(node), BLACK)
				setColor(rightOf(brother), BLACK)
				tree.rotateLeft(parentOf(node))

				node = tree.root
			}
		} else {
			// 要删除的节点在父亲右边，对应图例3，4
			// 找出兄弟
			brother := leftOf(parentOf(node))

			// 兄弟是红色的，对应图例3，那么兄弟变黑，父亲变红，然后对父亲右旋，进入图例41,42,43
			if isRed(brother) {
				setColor(brother, BLACK)
				setColor(parentOf(node), RED)
				tree.rotateRight(parentOf(node))
				brother = leftOf(parentOf(node)) // 图例3调整后进入图例41,42,43，兄弟此时变了
			}

			// 兄弟是黑色的，对应图例41，42，43
			// 兄弟的左右儿子都是黑色，进入图例43，将兄弟设为红色，父亲所在的子树作为整体，当作删除的节点，继续向上递归
			if !isRed(leftOf(brother)) && !isRed(rightOf(brother)) {
				setColor(brother, RED)
				node = parentOf(node)
			} else {
				// 兄弟的左儿子是黑色，进入图例42，将兄弟设为红色，兄弟的右儿子设为黑色，对兄弟左旋，进入图例41
				if !isRed(leftOf(brother)) {
					setColor(rightOf(brother), BLACK)
					setColor(brother, RED)
					tree.rotateLeft(brother)
					brother = leftOf(parentOf(node)) // 图例42调整后进入图例41，兄弟此时变了
				}

				// 兄弟的左儿子是红色，进入图例41，将兄弟设置为父亲的颜色，兄弟的左儿子以及父亲变黑，对父亲右旋
				setColor(brother, parentOf(node).color)
				setColor(parentOf(node), BLACK)
				setColor(leftOf(brother), BLACK)
				tree.rotateRight(parentOf(node))

				node = tree.root
			}
		}
	}

	// this node always black
	setColor(node, BLACK)
}

// find min key pairs
func (tree *rbTree) MinKey() (key string, value interface{}, exist bool) {
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

func (node *rbTNode) minNode() *rbTNode {
	// 左子树为空，表面已经是最左的节点了，该值就是最小值
	if node.left == nil {
		return node
	}

	// 一直左子树递归
	return node.left.minNode()
}

// find max key pairs
func (tree *rbTree) MaxKey() (key string, value interface{}, exist bool) {
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

func (node *rbTNode) maxNode() *rbTNode {
	// 右子树为空，表面已经是最右的节点了，该值就是最大值
	if node.right == nil {
		return node
	}

	// 一直右子树递归
	return node.right.maxNode()
}

// 查找指定节点
func (tree *rbTree) Get(key string) (value interface{}, exist bool) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	if tree.root == nil {
		// 如果是空树，返回空
		return
	}

	node := tree.root.find(key)
	if node == nil {
		return nil, false
	}
	return node.v, true
}

// 查找指定节点
func (tree *rbTree) Contains(key string) (exist bool) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	if tree.root == nil {
		// 如果是空树，返回空
		return
	}
	node := tree.root.find(key)
	if node == nil {
		return false
	}
	return true
}

func (tree *rbTree) Len() int64 {
	return tree.len
}

// get int
func (tree *rbTree) GetInt(key string) (value int, exist bool, err error) {
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
func (tree *rbTree) GetInt64(key string) (value int64, exist bool, err error) {
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
func (tree *rbTree) GetString(key string) (value string, exist bool, err error) {
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
func (tree *rbTree) GetFloat64(key string) (value float64, exist bool, err error) {
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
func (tree *rbTree) GetBytes(key string) (value []byte, exist bool, err error) {
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

// 查找指定节点
func (tree *rbTree) find(key string) *rbTNode {
	if tree.root == nil {
		// 如果是空树，返回空
		return nil
	}

	return tree.root.find(key)
}

func (node *rbTNode) find(key string) *rbTNode {
	cmp := compare(key, node.k)
	if cmp == 0 {
		// 如果该节点刚刚等于该值，那么返回该节点
		return node
	} else if cmp < 0 {
		// 如果查找的值小于节点值，从节点的左子树开始找
		if node.left == nil {
			// 左子树为空，表示找不到该值了，返回nil
			return nil
		}
		return node.left.find(key)
	} else {
		// 如果查找的值大于节点值，从节点的右子树开始找
		if node.right == nil {
			// 右子树为空，表示找不到该值了，返回nil
			return nil
		}
		return node.right.find(key)
	}
}

// 中序遍历
// mid order get key list
func (tree *rbTree) KeySortedList() []string {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	keyList := make([]string, 0, tree.len)
	return tree.root.midOrder(keyList)
}

func (node *rbTNode) midOrder(keyList []string) []string {
	if node == nil {
		return keyList
	}

	// 先打印左子树
	keyList = node.left.midOrder(keyList)

	keyList = append(keyList, node.k)

	// 打印右子树
	keyList = node.right.midOrder(keyList)

	return keyList
}

// 验证是不是棵红黑树
func (tree *rbTree) Check() bool {
	if tree == nil || tree.root == nil {
		return true
	}

	// 判断树是否是一棵二分查找树
	if !tree.root.isBST() {
		fmt.Println("is not BST")
		return false
	}

	// 判断树是否遵循2-3-4树，也就是不能有连续的两个红链接
	if !tree.root.is234() {
		fmt.Println("is not 234 tree")
		return false
	}

	// 判断树是否平衡，也就是任意一个节点到叶子节点，经过的黑色链接数量相同
	// 先计算根节点到最左边叶子节点的黑链接数量
	blackNum := 0
	x := tree.root
	for x != nil {
		if !isRed(x) { // 是黑色链接
			blackNum = blackNum + 1
		}
		x = x.left
	}

	if !tree.root.isBalanced(blackNum) {
		fmt.Println("is not Balanced")
		return false
	}
	return true
}

// 节点所在的子树是否是一棵二分查找树
func (node *rbTNode) isBST() bool {
	if node == nil {
		return true
	}

	// 左子树非空，那么根节点必须大于左儿子节点
	if node.left != nil {
		cmp := compare(node.k, node.left.k)
		if cmp > 0 {
		} else {
			fmt.Printf("father:%#v,lchild:%#v,rchild:%#v\n", node, node.left, node.right)
			return false
		}
	}

	// 右子树非空，那么根节点必须小于右儿子节点
	if node.right != nil {
		cmp := compare(node.k, node.right.k)
		if cmp < 0 {
		} else {
			fmt.Printf("father:%#v,lchild:%#v,rchild:%#v\n", node, node.left, node.right)
			return false
		}
	}

	// 左子树也要判断是否是平衡查找树
	if !node.left.isBST() {
		return false
	}

	// 右子树也要判断是否是平衡查找树
	if !node.right.isBST() {
		return false
	}

	return true
}

// 节点所在的子树是否遵循2-3-4树
func (node *rbTNode) is234() bool {
	if node == nil {
		return true
	}

	// 不允许连续两个左红链接
	if isRed(node) && isRed(node.left) {
		fmt.Printf("father:%#v,lchild:%#v\n", node, node.left)
		return false
	}

	if isRed(node) && isRed(node.right) {
		fmt.Printf("father:%#v,rchild:%#v\n", node, node.right)
		return false
	}

	// 左子树也要判断是否遵循2-3-4树
	if !node.left.is234() {
		return false
	}

	// 右子树也要判断是否是遵循2-3-4树
	if !node.right.is234() {
		return false
	}

	return true
}

// 节点所在的子树是否平衡，是否有 blackNum 个黑链接
func (node *rbTNode) isBalanced(blackNum int) bool {
	if node == nil {
		return blackNum == 0
	}

	if !isRed(node) {
		blackNum = blackNum - 1
	}

	if !node.left.isBalanced(blackNum) {
		fmt.Println("node.left to leaf black link is not ", blackNum)
		return false
	}

	if !node.right.isBalanced(blackNum) {
		fmt.Println("node.right to leaf black link is not ", blackNum)
		return false
	}

	return true
}

// iterator help struct
type bsTreeNode interface {
	leftOf() bsTreeNode
	rightOf() bsTreeNode
	values() (key string, value interface{})
}

// 返回节点的左子节点
func (node *rbTNode) leftOf() bsTreeNode {
	if node.left == nil {
		return nil
	}

	return node.left
}

// 返回节点的右子节点
func (node *rbTNode) rightOf() bsTreeNode {
	if node.right == nil {
		return nil
	}

	return node.right
}

// not check node nil, may be panic, user should deal by oneself
func (node *rbTNode) values() (key string, value interface{}) {
	return node.k, node.v
}

func (tree *rbTree) KeyList() []string {
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

func (tree *rbTree) Iterator() MapIterator {
	q := new(linkQueue)
	if tree.root != nil {
		q.add(tree.root)
	}
	return q
}
