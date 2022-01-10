Miscellaneous rope tests during development

/*
func main() {
	// testing block for rotation left
	/*
		t1 := &AVLNode{Value: []rune{'a'}, Height: 0, Weight: 1}
		t23 := &AVLNode{Value: []rune{'b', 'c'}, Height: 0, Weight: 2}
		t4 := &AVLNode{Value: []rune{'d'}, Height: 0, Weight: 1}

		z := mkNode(nil, t23, t4)
		x := mkNode(nil, t1, z)


		p := x.rotateLeft()
		PrintInline(p)
		LeavesOnly(p)

		TreeStats(p)
		fmt.Println(x.Weight)
		fmt.Println(z.Weight)

		// testing block for rotation right

		p = p.rotateRight()
		PrintInline(p)
		LeavesOnly(p)

		TreeStats(p)
		fmt.Println(x.Weight)
		fmt.Println(z.Weight)
	*/

	// testing basic join
	/*
		a := &AVLNode{Value: []rune{'a', 'b', 'c'}, Height: 0, Weight: 3}
		b := &AVLNode{Value: []rune{'d', 'e'}, Height: 0, Weight: 2}
		c := &AVLNode{Value: []rune{'f', 'g'}, Height: 0, Weight: 2}
		d := &AVLNode{Value: []rune{'h', 'i'}, Height: 0, Weight: 2}

		n1 := mkNode(nil, b, c)
		n2 := mkNode(nil, a, n1)

		z := Join(n2, d)
		z.updateWeight()
		e := &AVLNode{L: z}
		z.U = e

		e.updateHeight()
		e.Weight = e.L.sumTree()
		fmt.Println(e.L.sumTree())

		buf := &bytes.Buffer{}
		memviz.Map(buf, e)
		err := ioutil.WriteFile("tree-data", buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	*/
	// testing normal creation
	/*
		file, err := os.Open("test.txt")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		// The file.Read() function will happily read a tiny file in to a large
		// byte slice, but io.ReadFull() will return an
		// error if the file is smaller than the byte slice.
		byteSlice := make([]byte, 90)
		_, err = io.ReadFull(file, byteSlice)
		if err != nil {
			log.Fatal(err)
		}
		r := mkRope()
		r.Load(byteSlice)
		PrintInline(r.Head)
		LeavesOnly(r.Head)
		TreeStats(r.Head)
		CheckAVL(r.Head)
		buf := &bytes.Buffer{}
		memviz.Map(buf, r)
		err = ioutil.WriteFile("tree-data", buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	*/

	// testing split
	/*
		a := &AVLNode{Value: []rune{'a', 'b', 'c'}, Height: 0, Weight: 3}
		b := &AVLNode{Value: []rune{'d', 'e'}, Height: 0, Weight: 2}
		c := &AVLNode{Value: []rune{'f', 'g'}, Height: 0, Weight: 2}
		d := &AVLNode{Value: []rune{'h', 'i'}, Height: 0, Weight: 2}

		n1 := mkNode(nil, b, c)
		n2 := mkNode(nil, a, n1)

		z := Join(n2, d)
		buf := &bytes.Buffer{}
		memviz.Map(buf, z)
		err := ioutil.WriteFile("tree-data", buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
		l, r := Split(z, 3)
		PrintInline(l)
		PrintInline(r)
		fmt.Println(l.Weight)
		fmt.Println(r.Weight)
		TreeStats(r)

		CheckAVL(l)
		CheckAVL(r)
	*/

	/*
		// testing ranged reports
		file, err := os.Open("test.txt")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		// The file.Read() function will happily read a tiny file in to a large
		// byte slice, but io.ReadFull() will return an
		// error if the file is smaller than the byte slice.
		byteSlice := make([]byte, 90)
		_, err = io.ReadFull(file, byteSlice)
		if err != nil {
			log.Fatal(err)
		}
		r := mkRope()
		r.Load(byteSlice)
		res, err := r.ReportRange(0, 11)
		if err != nil {
			return
		}
		PrintInline(r.Head)
		fmt.Printf("/%s/", string(res))
	*/

	g := mkGapBuf()

	file, err := os.Open("testing.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ret := make([]byte, 200)
	_, err = io.ReadFull(file, ret)
	if err != nil {
		panic(err)
	}

	app := make([]rune, len(ret))
	for i, bt := range ret {
		app[i] = rune(bt)
	}

	g.Load(ret[:100])
	g.Append(app)
	g.Load(ret)

	//fmt.Println(string(g.GapStart.Next.Content))
	//fmt.Println(g.Content.Head.Next.Content)

	fmt.Println(g.ToString())
	res, err := g.Report()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))

	g.Append(app)
	fmt.Println()
	fmt.Println(g.ToString())
	res, err = g.Report()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

}