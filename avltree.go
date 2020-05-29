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

// Deprecate
// AVL Tree
// Use recursion.
type avlTree struct {
	c          comparator   // tree key compare
	root       *avlTreeNode // tree root node
	len        int64        // tree key pairs num
	sync.Mutex              // lock for concurrent safe
}

// AVL节点
type avlTreeNode struct {
	k      string       // key
	v      interface{}  // value
	height int64        // 该节点作为树根节点，树的高度，方便计算平衡因子
	left   *avlTreeNode // 左子树
	right  *avlTreeNode // 右字树
}

func (node *avlTreeNode) h() int64 {
	if node == nil {
		return 0
	}

	lh := node.left.h()
	rh := node.right.h()
	if lh > rh {
		return lh + 1
	} else {
		return rh + 1
	}
}

func (tree *avlTree) Height() int64 {
	return tree.root.h()
}

// 更新节点的树高度
func (node *avlTreeNode) updateHeight() {
	if node == nil {
		return
	}

	var leftHeight, rightHeight int64 = 0, 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}

	//leftHeight, rightHeight := height(node.left), height(node.right)

	// 哪个子树高算哪棵的
	maxHeight := leftHeight
	if rightHeight > maxHeight {
		maxHeight = rightHeight
	}
	// 高度加上自己那一层
	node.height = maxHeight + 1
}

// 计算平衡因子
func (node *avlTreeNode) balanceFactor() int64 {
	//return height(node.left) - height(node.right)
	var leftHeight, rightHeight int64 = 0, 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	return leftHeight - rightHeight
}

// 单右旋操作，看图说话
func (_ *avlTreeNode) rightRotation(Root *avlTreeNode) *avlTreeNode {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%#v\n,l:%#v\n,r:%#v\n", Root, Root.left, Root.right)
			panic(err)
		}
	}()
	// 只有Pivot和B，Root位置变了
	Pivot := Root.left
	B := Pivot.right
	Pivot.right = Root
	Root.left = B

	// 只有Root和Pivot变化了高度
	Root.updateHeight()
	Pivot.updateHeight()
	return Pivot
}

// 单左旋操作，看图说话
func (_ *avlTreeNode) leftRotation(Root *avlTreeNode) *avlTreeNode {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%#v\n,l:%#v\n,r:%#v\n", Root, Root.left, Root.right)
			panic(err)
		}
	}()
	// 只有Pivot和B，Root位置变了
	Pivot := Root.right
	B := Pivot.left
	Pivot.left = Root
	Root.right = B

	// 只有Root和Pivot变化了高度
	Root.updateHeight()
	Pivot.updateHeight()
	return Pivot
}

// 先左后右旋操作，看图说话
func (_ *avlTreeNode) leftRightRotation(node *avlTreeNode) *avlTreeNode {
	node.left = node.leftRotation(node.left)
	return node.rightRotation(node)
}

// 先右后左旋操作，看图说话
func (_ *avlTreeNode) rightLeftRotation(node *avlTreeNode) *avlTreeNode {
	node.right = node.rightRotation(node.right)
	return node.leftRotation(node)
}

// set Comparator
func (tree *avlTree) SetComparator(c comparator) {
	if tree.len == 0 {
		tree.c = c
	}
}

// 添加元素
func (tree *avlTree) Put(key string, value interface{}) {
	// add lock
	tree.Lock()
	defer tree.Unlock()

	add := false
	if tree.root != nil {
		node := tree.root.find(tree.c, key)
		if node == nil {
			add = true
		}
	} else {
		add = true
	}

	// 往树根添加元素，会返回新的树根
	tree.root = tree.root.put(tree.c, key, value)

	if add {
		tree.len = tree.len + 1
	}
}

func (node *avlTreeNode) put(compare comparator, key string, value interface{}) *avlTreeNode {
	// 添加值到根节点node，如果node为空，那么让值成为新的根节点，树的高度为1
	if node == nil {
		return &avlTreeNode{k: key, v: value, height: 1}
	}

	// 如果值重复，什么都不用做，直接更新次数
	if node.k == key {
		node.v = value
		return node
	}

	// 辅助变量
	var newTreeNode *avlTreeNode

	cmp := compare(key, node.k)
	if cmp > 0 {
		// 插入的值大于节点值，要从右子树继续插入
		node.right = node.right.put(compare, key, value)
		// 平衡因子，插入右子树后，要确保树根左子树的高度不能比右子树低一层。
		factor := node.balanceFactor()
		// 右子树的高度变高了，导致左子树-右子树的高度从-1变成了-2。
		if factor == -2 {
			cmp1 := compare(key, node.right.k)
			if cmp1 > 0 {
				// 表示在右子树上插上右儿子导致失衡，需要单左旋：
				newTreeNode = node.leftRotation(node)
			} else {
				//表示在右子树上插上左儿子导致失衡，先右后左旋：
				newTreeNode = node.rightLeftRotation(node)
			}
		}
	} else {
		// 插入的值小于节点值，要从左子树继续插入
		node.left = node.left.put(compare, key, value)
		// 平衡因子，插入左子树后，要确保树根左子树的高度不能比右子树高一层。
		factor := node.balanceFactor()
		// 左子树的高度变高了，导致左子树-右子树的高度从1变成了2。
		if factor == 2 {
			cmp1 := compare(key, node.left.k)
			if cmp1 < 0 {
				// 表示在左子树上插上左儿子导致失衡，需要单右旋：
				newTreeNode = node.rightRotation(node)
			} else {
				//表示在左子树上插上右儿子导致失衡，先左后右旋：
				newTreeNode = node.leftRightRotation(node)
			}
		}
	}

	if newTreeNode == nil {
		// 表示什么旋转都没有，根节点没变，直接刷新树高度
		node.updateHeight()
		return node
	} else {
		// 旋转了，树根节点变了，需要刷新新的树根高度
		newTreeNode.updateHeight()
		return newTreeNode
	}
}

// find min key pairs
func (tree *avlTree) MinKey() (key string, value interface{}, exist bool) {
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

func (node *avlTreeNode) minNode() *avlTreeNode {
	// 左子树为空，表面已经是最左的节点了，该值就是最小值
	if node.left == nil {
		return node
	}

	// 一直左子树递归
	return node.left.minNode()
}

// 找出最大值的节点
func (tree *avlTree) MaxKey() (key string, value interface{}, exist bool) {
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

func (node *avlTreeNode) maxNode() *avlTreeNode {
	// 右子树为空，表面已经是最右的节点了，该值就是最大值
	if node.right == nil {
		return node
	}

	// 一直右子树递归
	return node.right.maxNode()
}

// 查找指定节点
func (tree *avlTree) Get(key string) (value interface{}, exist bool) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	if tree.root == nil {
		// 如果是空树，返回空
		return
	}

	node := tree.root.find(tree.c, key)
	if node == nil {
		return nil, false
	}
	return node.v, true
}

func (node *avlTreeNode) find(compare comparator, key string) *avlTreeNode {
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
		return node.left.find(compare, key)
	} else {
		// 如果查找的值大于节点值，从节点的右子树开始找
		if node.right == nil {
			// 右子树为空，表示找不到该值了，返回nil
			return nil
		}
		return node.right.find(compare, key)
	}
}

// 删除指定的元素
func (tree *avlTree) Delete(key string) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	if tree.len == 0 {
		return
	}

	// 查找元素是否存在，不存在则退出
	node := tree.root.find(tree.c, key)
	if node == nil {
		return
	}

	tree.root = tree.root.delete(tree.c, key)
	tree.len = tree.len - 1

}

func (node *avlTreeNode) delete(compare comparator, key string) *avlTreeNode {
	if node == nil {
		// 如果是空树，直接返回
		return nil
	}

	cmp := compare(key, node.k)
	if cmp < 0 {
		// 从左子树开始删除
		node.left = node.left.delete(compare, key)
		// 删除后要更新该子树高度
		node.left.updateHeight()
	} else if cmp > 0 {
		// 从右子树开始删除
		node.right = node.right.delete(compare, key)
		// 删除后要更新该子树高度
		node.right.updateHeight()
	} else {
		// 找到该值对应的节点
		// 该节点没有左右子树
		// 第一种情况，删除的节点没有儿子，直接删除即可。
		if node.left == nil && node.right == nil {
			return nil // 直接返回nil，表示直接该值删除
		}

		// 该节点有两棵子树，选择更高的哪个来替换
		// 第二种情况，删除的节点下有两个子树，选择高度更高的子树下的节点来替换被删除的节点，如果左子树更高，选择左子树中最大的节点，也就是左子树最右边的叶子节点，如果右子树更高，选择右子树中最小的节点，也就是右子树最左边的叶子节点。最后，删除这个叶子节点。
		if node.left != nil && node.right != nil {
			// 左子树更高，拿左子树中最大值的节点替换
			if node.left.height > node.right.height {
				maxNode := node.left
				for maxNode.right != nil {
					maxNode = maxNode.right
				}

				// 最大值的节点替换被删除节点
				node.k = maxNode.k
				node.v = maxNode.v

				// 把最大的节点删掉
				node.left = node.left.delete(compare, maxNode.k)
				// 删除后要更新该子树高度
				node.left.updateHeight()
			} else {
				// 右子树更高，拿右子树中最小值的节点替换
				minNode := node.right
				for minNode.left != nil {
					minNode = minNode.left
				}

				// 最小值的节点替换被删除节点
				node.k = minNode.k
				node.v = minNode.v

				// 把最小的节点删掉
				node.right = node.right.delete(compare, minNode.k)
				// 删除后要更新该子树高度
				node.right.updateHeight()
			}
		} else {
			// 只有左子树或只有右子树
			// 只有一个子树，该子树也只是一个节点，将该节点替换被删除的节点，然后置子树为空
			if node.left != nil {
				//第三种情况，删除的节点只有左子树，因为树的特征，可以知道左子树其实就只有一个节点，它本身，否则高度差就等于2了。
				node.k = node.left.k
				node.v = node.left.v
				node.height = 1
				node.left = nil
			} else if node.right != nil {
				//第四种情况，删除的节点只有右子树，因为树的特征，可以知道右子树其实就只有一个节点，它本身，否则高度差就等于2了。
				node.k = node.right.k
				node.v = node.right.v
				node.height = 1
				node.right = nil
			}
		}

		// 找到值后，进行替换删除后，直接返回该节点
		return node
	}

	// 左右子树递归删除节点后需要平衡
	var newNode *avlTreeNode
	// 相当删除了右子树的节点，左边比右边高了，不平衡
	if cmp > 0 && node.balanceFactor() == 2 {
		//fmt.Println("l-r=2 and l:", node.left.balanceFactor())
		if node.left.balanceFactor() >= 0 { // why >0 will err must be checking
			newNode = node.rightRotation(node)
		} else {
			newNode = node.leftRightRotation(node)
		}
		//  相当删除了左子树的节点，右边比左边高了，不平衡
	} else if cmp < 0 && node.balanceFactor() == -2 {
		//fmt.Println("l-r=-2 and l:", node.right.balanceFactor())
		if node.right.balanceFactor() <= 0 {
			newNode = node.leftRotation(node)
		} else {
			newNode = node.rightLeftRotation(node)
		}
	}

	if newNode == nil {
		node.updateHeight()
		return node
	} else {
		newNode.updateHeight()
		return newNode
	}
}

// 查找指定节点
func (tree *avlTree) Contains(key string) (exist bool) {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	if tree.root == nil {
		// 如果是空树，返回空
		return
	}
	node := tree.root.find(tree.c, key)
	if node == nil {
		return false
	}
	return true
}

func (tree *avlTree) Len() int64 {
	return tree.len
}

// get int
func (tree *avlTree) GetInt(key string) (value int, exist bool, err error) {
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
func (tree *avlTree) GetInt64(key string) (value int64, exist bool, err error) {
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
func (tree *avlTree) GetString(key string) (value string, exist bool, err error) {
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
func (tree *avlTree) GetFloat64(key string) (value float64, exist bool, err error) {
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
func (tree *avlTree) GetBytes(key string) (value []byte, exist bool, err error) {
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

// 中序遍历
// mid order get key list
func (tree *avlTree) KeySortedList() []string {
	// add lock
	tree.Lock()
	defer tree.Unlock()
	keyList := make([]string, 0, tree.len)
	return tree.root.midOrder(keyList)
}

func (node *avlTreeNode) midOrder(keyList []string) []string {
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

// 验证是不是棵AVL树
func (tree *avlTree) Check() bool {
	if tree == nil || tree.root == nil {
		return true
	}

	// 判断节点是否符合 AVL 树的定义
	if tree.root.isAVL(tree.c) {
		return true
	}

	return false
}

// 判断节点是否符合 AVL 树的定义
func (node *avlTreeNode) isAVL(compare comparator) bool {
	if node == nil {
		return true
	}

	// 左右子树都为空，那么是叶子节点
	if node.left == nil && node.right == nil {
		// 叶子节点高度应该为1
		if node.height == 1 {
			return true
		} else {
			fmt.Println("leaf node height must is 1, now is:", node.height)
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

		bal := node.left.height - node.right.height
		if bal < 0 {
			bal = -bal
		}

		// 子树高度差不能大于1
		if bal > 1 {
			fmt.Println("sub tree height bal is ", bal)
			return false
		}

		// 如果左子树比右子树高，那么父亲的高度等于左子树+1
		if node.left.height > node.right.height {
			if node.height == node.left.height+1 {
			} else {
				fmt.Printf("%#v height:%v,left sub tree height: %v,right sub tree height:%v\n", node, node.height, node.left.height, node.right.height)
				return false
			}
		} else {
			// 如果右子树比左子树高，那么父亲的高度等于右子树+1
			if node.height == node.right.height+1 {
			} else {
				fmt.Printf("%#v height:%v,left sub tree height: %v,right sub tree height:%v\n", node, node.height, node.left.height, node.right.height)
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
			if node.right.height == 1 && node.right.left == nil && node.right.right == nil {
				cmp := compare(node.right.k, node.k)
				if cmp > 0 {
					// 右节点必须比父亲大
				} else {
					fmt.Printf("%v has only right tree,but right small", node.k)
					return false
				}
			} else {
				fmt.Printf("%v has only right tree height:%d,but right has lc: %#v, rc：%#v", node.k, node.right.height, node.right.left, node.right.right)
				return false
			}
		} else {
			if node.left.height == 1 && node.left.left == nil && node.left.right == nil {
				cmp := compare(node.left.k, node.k)
				if cmp < 0 {
					// 左节点必须比父亲小
				} else {
					fmt.Printf("%v has only left tree,but right small", node.k)
					return false
				}
			} else {
				fmt.Printf("%v has only left tree height:%d,but left has lc: %#v, rc：%#v", node.k, node.left.height, node.left.left, node.left.right)
				return false
			}
		}
	}

	return true
}

// 返回节点的左子节点
func (node *avlTreeNode) leftOf() bsTreeNode {
	if node.left == nil {
		return nil
	}

	return node.left
}

// 返回节点的右子节点
func (node *avlTreeNode) rightOf() bsTreeNode {
	if node.right == nil {
		return nil
	}

	return node.right
}

// not check node nil, may be panic, user should deal by oneself
func (node *avlTreeNode) values() (key string, value interface{}) {
	return node.k, node.v
}

func (tree *avlTree) KeyList() []string {
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

func (tree *avlTree) Iterator() MapIterator {
	q := new(linkQueue)
	if tree.root != nil {
		q.add(tree.root)
	}
	return q
}
