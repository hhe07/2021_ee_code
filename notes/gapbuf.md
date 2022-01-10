```go
type StorageType interface {
	Report() ([]rune, error)              // report entire buffer
	ReportRange(i, j int) ([]rune, error) // report segment of buffer
	ReportCharacter(i int) (rune, error)  // report single character
	Insert(i int, content []rune) error
	Append(i int, content []rune) error
	Replace(i int, content []rune) error
	Split(i int) (StorageType, error) // TODO: will this accept both rope and gapbuffer?
	DeleteRange(i, j int) ([]rune, error)
	Concat(content []rune) error
	Save(f *FileWrapper) error
    IsReady() (bool, error)
	Load(f *FileWrapper) error
}

```

Gap Buffer Implementation Notes:


## getNode (done)
Returns the node at index ``i``

Implementation:
- check index
- if index is less than gap start, do normal
- else if index is greater than gap start, subtract gap length
- return value of get from linked list

Testing:
- test start and end for their respective indices
- test whether out of bound indices throw an error (it should)
- test before and after gap (somewhere in the middle)


## indexCheck (done)
- check empty
- check in bounds
- check that in order from least to greatest

## checkEmpty (done)
- check that length of list is not zero


## Report (done)
A special case of the ``ReportRange`` method (see below), in which indices selected
are ``[0, len(list)-1]`` 

## ReportRange (done)
Returns the range of characters in ``[i,j]``.

Implementation:
- Check indices for no errors (write indexCheck method)
- Create slice to hold results
- Get node i (``getNode``)
- For node between i and j (inclusive)
  - If node is not in gap:
    - add content to slice
  - move node to next node
- return result

Testing:
- Test start and end for their respective indices
- Test whether j < i throws an error (it should)
- Test whether i = j throws an error (it should)
- test whether indices out of bound throw an error (it should)
- Test somewhere in the middle



## ReportCharacter (done)
Returns the character at index.

Implementation:
- check index
- get node
- return value of node

Testing:
- test invalid indices
- test around gap / "in" gap to make sure that the gap isn't being counted
- test after some gap movements to both make sure that the gap is properly moved and to make sure that an index inside the gap isn't counted
- test at start and end

## insertAtGap (done)
Inserts the slice ``content`` at the current gap position

Implementation:
- get gap start
- if length of gap will be exceeded:
  - expand gap (``expandGap``)
- for item in ``content``:
  - fill current empty node with content
  - advance gap start
  - advance current node
  - decrease gap length and increase gap start idx

Testing:
- test that the gap is properly expanded if overflowed (right size: should always be 2 larger)
- honestly will be tested quite a bit in Insert / Append.

## expandGap (done)
Expands the gap at the current place.

Implementation:
- Delink and save the pointer of the next actually storing item after the gap
- For the extra size needed:
  - Create a new link and link it to the current end
    - TODO: use the actual function under the ll class, so you don't need to increase internal list length in GapBuffer
  - Change the gap end position to its next (should always be the next link)
  - Increase gap length
  - Increase length of internal LinkedList content

Testing:
- Test various lengths

## Insert (done)

Inserts a slice ``content`` of rune at a certain index i.

Implementation:
- Move gap to right location
- insertAtGap ``content``


## Append
Inserts a slice ``content`` at the end of the buffer
Implementation:
- Move gap to end
- insertAtGap ``content``

## moveGap
Moves the gap such that it begins at position ``i``.

If i == length, the gap is at the very end of the buffer (no characters after).

Implementation (current):
- if i == current position, exit
- check index
- Get node at (i-1): this must be done beforehand to prevent wrong index from being used.
- if gap is currently at start:
  - set list start to element after gap
  - only delink end and element after
- if gap is currently at end:
  - set list end to element before gap
  - only delink element before
- else:
  - Store elements before head and after end
  - Delink head and end from both.
  - Link before head and after end together
- if i == 0:
  - only link gap end to current head
  - set head to gap start
- if i == length:
  - only link gap start to current end
  - set end to gap end
- else:
  - Delink and store node after (i-1)
  - Link head to (i-1)
  - Link end to node after (i-1)
- Done

Testing:
- bad index
- move gap to start, end
- move gap to somewhere in the middle.

## Replace (done)
replace the content starting from index i.

Implementation:
- check index
- get node i
- for rune in content:
  - replace content in current node with rune
  - get next node
  - if the next is nil, return an error
  - if the next is in the buffer, enter a loop to skip over.

Testing:
- bad index
- replacement at start, near the end
- replacement at centre
- overflowing replacement
- making sure it flows around the gap.

## Split
Split the current buffer such that two results are produced:
[0, i] and [i+1, len-1].

The new buffer has the gap initialised at the start.
Splitting MUST be at a point in the centre

Implementation:
- check index
- Move gap of current buffer back to start, compress down to 0
- Split gap at point
- Remake gap in current
- Make gap in new
- Set lengths right
- return new

Test:
- Invalid index
- Split at start / end (shouldn't work)

## compressGap
Compresses current gap to size 0.

Implementation
- If gap is already length 0:
  - exit
- If gap at start:
  - make pointer to first the element beyond gap
  - delink from gap buffer
- if gap at end
  - make pointer to end the first element before gap
  - delink from gap buffer
- else:
  - delink gap
  - Link before and after gap
- reduce length of internal list
- remove pointers to either side of gap
- reduce gap size and index to 0

Test:
- just do the thing

## deleteRange
Extend gap over the additional indices

Implementation
- Check indices
- Get number of nodes that need to be incorporated into the gap
- Move gap to position
- For number of nodes > 0:
  - move gap end 1 position
  - add current value to list
  - extend gap by 1

Test:
- overflowing indices
- Delete from start and end
- delete in a central position

## Concat
Append ``content`` onto the list

Implementation:
- Just duplicated append

Test:
- Covered in Append

## Save
TODO: work

## Load:
- Open file
- Create internal list
- Append Content
- make gap

Test:
- Load a file (obviously)

## ToString:
keep same

```go
func (g *GapBuffer) CursorTo(i int) error {
	// todo: work on
	if g.GapStartIdx == i {
		return nil
	}
	// special case: if i = 0
	if i == 0 {
		tmp := g.Content.Head
		// couple prior and after gap together
		prior := g.GapStart.Prior
		end := g.GapEnd.Next
		prior.Next = end
		end.Prior = prior

		g.GapStart.Prior = nil

		g.Content.Head = g.GapStart
		g.GapEnd.Next = tmp
		tmp.Prior = g.Content.Head
		g.GapStartIdx = 0
		return nil
	} else if i >= g.Length() {
		prior := g.GapStart.Prior
		end := g.GapEnd.Next
		prior.Next = end
		end.Prior = prior

		tmp := g.Content.End
		g.GapEnd.Next = nil
		g.Content.End = g.GapEnd
		g.GapStart.Prior = tmp
		tmp.Next = g.GapStart
		g.GapStartIdx = g.Length()
		return nil
	}
	tAmt := i - g.GapStartIdx
	if tAmt > 0 {
		// move gap forwards
		for tAmt > 0 {
			mv := g.GapEnd.Next.Content // moved character
			g.GapStart.Content = mv
			g.GapStart = g.GapStart.Next
			g.GapEnd = g.GapEnd.Next
			g.GapEnd.Content = -1

			tAmt--
			g.GapStartIdx++
		}
		return nil
	}
	// move gap backwards
	for tAmt < 0 {
		mv := g.GapStart.Prior.Content // moved character
		g.GapEnd.Content = mv
		g.GapEnd = g.GapEnd.Prior
		g.GapStart = g.GapStart.Prior
		g.GapStart.Content = -1

		tAmt++
		g.GapStartIdx--
	}
	return nil

	/*
		TODO: this is good but doesnt work

		// moves the gap as a block to position i
		el, err := g.getNode(i - 1)
		if err != nil {
			return err
		}
		el.Next.Prior = g.GapEnd
		g.GapEnd.Next = el.Next

		el.Next = g.GapStart
		g.GapStart.Prior = el

		return nil

		// displaces
	*/
	// todo: shortcut for when i is already start index, start and end positions
	// todo: define behaviour well: does it displace indicated idx or otherwise

}

func (g *GapBuffer) DeleteRange(i, j int) ([]rune, error) {
	// TODO: single character removal doesn't account for gap
	if i == j {
		r, err := g.Content.Remove(i)
		return []rune{r}, err
	}
	delStart, err := g.getNode(i) // included: at position i-1
	if err != nil {
		return nil, err
	}
	delEnd, err := g.getNode(j) // excluded
	if err != nil {
		return nil, err
	}
	// case: start
	if i == 0 {
		g.Content.Head = delEnd
		delEnd.Prior = nil
	} else if j == g.Length()-1 {
		// todo: fix
		g.Content.End = delStart
		g.Content.End.Next = nil
	} else {
		delStart.Prior.Next = delEnd
		delEnd.Prior = delStart.Prior
	}
	ret := make([]rune, 0)
	for delStart != nil {
		ret = append(ret, delStart.Content)
		delStart = delStart.Next
		g.Content.Length--
	}
	return ret, nil
}


/*
func TestInsertA(t *testing.T) {
	t.Log("open file")
	f := FileWrapper{
		Filename: "testlist.md",
		CharLen:  0,
	}
	gb := GapBuffer{}
	gb.Load(&f)
	t.Log(gb.ToString())

	t.Log(gb.GapLen)

	t.Log()
	t.Log()
	err := gb.moveGap(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	t.Log()
	t.Log()

	t.Log(gb.Length())
	t.Log(gb.GapLen)
	gb.insertAtGap([]rune{'g', 'a', 'p', ' ', 'b', 'u', 'f', 'f', 'e', 'r', ' ', 'a', 'n', 'd'})
	t.Log(gb.ToString())
	t.Log()
	t.Log()

	t.Log(gb.Length())
	gb.Append([]rune{'w', 'a', 'w'})
	t.Log(gb.ToString())

	gb.DeleteRange(0, 2)
	t.Log(gb.ToString())
	// todo: remove, remove block, replace, split

}
*/

```