# Blockchain Merkle Patricia Trie


## Project 2: Build a Private BlockChain

Data structures:

### A) Block:

Each block must contain a header, and in the header there are the following fields: 

(1) Height: int32
(2) Timestamp: int64
The value must be in the UNIX timestamp format such as 1550013938
(3) Hash: string.

Block’s hash is the SHA3-256 encoded value of this string(note that you have to follow this specific order): 

```golang
hash_str := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + b.Value.Root + string(b.Header.Size)
```

(4) ParentHash: string
(5) Size: int32
The size is the length of the byte array of the block value

Each block must have a value, which is a Merkle Patricia Trie. All the data are inserted in the MPT and then a block contains that MPT as the value. So the field definition is this: 
Value: mpt MerklePatriciaTrie

Here's the summary of block structure: 
Block{Header{Height, Timestamp, Hash, ParentHash, Size}, Value}

Required functions: 
If arguments or return type is not specified, feel free to define them yourself. You may change the function's name, but make a comment to indicate which function you are implementing.

#### 1. Initial()
Description: This function takes arguments(such as height, parentHash, and value of MPT type) and forms a block. This is a method of the block struct.

#### 2.  DecodeFromJson(jsonString)
Description: This function takes a string that represents the JSON value of a block as an input, and decodes the input string back to a block instance. Note that you have to reconstruct an MPT from the JSON string, and use that MPT as the block's value. 
Argument: a string of JSON format
Return value: a block instance

#### 3. EncodeToJSON()
Description: This function encodes a block instance into a JSON format string. Note that the block's value is an MPT, and you have to record all of the (key, value) pairs that have been inserted into the MPT in your JSON string. There's an example with details on Piazza. Here's a website that can encode and decode JSON string: Link (Links to an external site.)Links to an external site.
Argument: a block or you may define this as a method of the block struct
Return value: a string of JSON format

Example of a block's JSON(decoded from JSON string):

```JSON
{
    "hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
    "timeStamp":1234567890,
    "height":1,
    "parentHash":"genesis",
    "size":1174,
    "mpt":{
        "charles":"ge",
        "hello":"world"
    }
}
```

### B) BlockChain:

Each blockchain must contain two fields described below. Don't change the name or the data type. 
(1) Chain: map[int32][]Block
This is a map which maps a block height to a list of blocks. The value is a list so that it can handle the forks.
(2) Length: int32
Length equals to the highest block height.

Required functions:
If arguments or return type is not specified, feel free to define them yourself. You may change the function's name, but make a comment to indicate which function you are implementing.

#### 1. Get(height)
Description: This function takes a height as the argument, returns the list of blocks stored in that height or None if the height doesn't exist.
Argument: int32
Return type: []Block

#### 2. Insert(block)
Description: This function takes a block as the argument, use its height to find the corresponding list in blockchain's Chain map. If the list has already contained that block's hash, ignore it because we don't store duplicate blocks; if not, insert the block into the list. 
Argument: block

#### 3. EncodeToJSON(self)
Description: This function iterates over all the blocks, generate blocks' JsonString by the function you implemented previously, and return the list of those JsonStritgns. 
Return type: string

Example of a blockchain's JSON:

```JSON
[
    {
        "hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
        "timeStamp":1234567890,
        "height":1,
        "parentHash":"genesis",
        "size":1174,
        "mpt":{
            "hello":"world",
            "charles":"ge"
        }
    },
    {
        "hash":"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159",
        "timeStamp":1234567890,
        "height":2,
        "parentHash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
        "size":1231,
        "mpt":{
            "hello":"world",
            "charles":"ge"
        }
    }
]
```

#### 4. DecodeFromJSON(self, jsonString)
Description: This function is called upon a blockchain instance. It takes a blockchain JSON string as input, decodes the JSON string back to a list of block JSON strings, decodes each block JSON string back to a block instance, and inserts every block into the blockchain. 
Argument: self, string

## Project 1: Build a Merkle Patricia Tree

Ethereum uses a Merkle Patricia Tree (Links to an external site.)Links to an external site. to store the transaction data in a block. By organizing the transaction data in a Merkle Patricia Tree, any block with fraudulent transactions would not match the tree's root hash. Build your own implementation of a Merkle Patricia Trie, following the specifications at the Ethereum wiki. (Links to an external site.)Links to an external site. Be mindful that you will use this code in the subsequent projects. 

Starter Code:

Rust: [Rust](https://github.com/CharlesGe129/cs686_blockchain_P1_Rust_skeleton.git) 

Go: [Go](https://github.com/CharlesGe129/cs686_blockchain_P1_Go_skeleton.git) 
 

Helpful resources:

Rust: 

[Rust1](https://doc.rust-lang.org/rust-by-example/index.html) 

[Rust2](https://doc.rust-lang.org/std/vec/?search=) 

Go:

[Go1](https://tour.golang.org/welcome/1) 

[Go2](https://golang.org/pkg/strings/#HasPrefix) 

 

## Project 1 specification

For this project, implement a Merkle Patricia Trie according to this Link (Links to an external site.)Links to an external site. and instructor's lectures. The functions and examples are written in RUST. Skeleton code and some help functions will be provided. You should not change the skeleton code. If you have different ideas that require to change the skeleton code, see the TA *before* submission. Any modification to the skeleton code without the TA's approval will result in point deduction.

You have to choose between Go and Rust. Skeleton code would be provided for both language and please write and pass your own tests before submit the final version of code.

You need to implement five features of the Merkle Patricia Trie:

### 1. Get()
Description: Get() function takes a key as the argument, traverse down the Merkle Patricia Trie and find the value. If the key doesn't exist, it will return an empty string.
Arguments: key(String)
Return: the value of that key in String type
Rust function definition: get(&mut self, key: &str) -> String

### 2. Insert()
Description: Insert() function takes a pair of <key, value> as arguments. It will traverse down the Merkle Patricia Trie, find the right place to insert the value, and do the insertion. 
Arguments: key(String), value(String)
Return: None
Rust function definition: get(&mut self, key: &str) -> String

### 3. Delete()
Description: Delete() function takes a key as the argument, traverse down the Merkle Patricia Trie and find that key. If the key exists, delete the corresponding value and re-balance the trie if necessary; if the key doesn't exist, return "path_not_found".
Arguments: key(String)
Return: None
Rust function definition: delete(&mut self, key: &str)

### 4. compact_encode()
Description: This function takes an array of HEX value as the input, mark the Node type(such as Branch, Leaf, Extension), make sure the length is even, and convert it into array of ASCII number as the output. You may find the Python version in this Link (Links to an external site.)Links to an external site., but the return type is different!
Arguments: hex_array(array of u8)
Return: array of u8
Rust function definition: compact_encode(hex_array: Vec<u8>) -> Vec<u8>
Example: input=[1, 6, 1], encoded_array=[1, 1, 6, 1], output=[17, 97]

### 5. compact_decode()
Description: This function reverses the compact_encode() function. 
Arguments: hex_array(array of u8)
Return: array of u8
Rust function definition: compact_decode(encoded_arr: Vec<u8>) -> Vec<u8>
Example: input=[17, 97], output=[1, 6, 1]

## Other help functions:

1. fn hash_node(node: &Node) -> String
Description: This function takes a node as the input, hash the node and return the hashed string.

If you use Golang, please follow this link to install the SHA3-256 package: https://github.com/golang/crypto (Links to an external site.)Links to an external site.

If you use Rust, the package dependency is written into Cargo.toml, so directly build and run the project should be fine. 

## Classes specification: 
In this project, there are two pre-defined classes. Both classes are defined in skeleton code, feel free to implement any useful functions.

1. enum Node
This class represent a node of type Branch, Leaf, Extension, or Null.

2. struct MerklePatriciaTrie
This class represent a Merkle Patricia Trie. It has two variables: "db" and "root".
Variable "db" is a HashMap. The key of the HashMap is a Node's hash value. The value of the HashMap is the Node. 
Variable "root" is a String, which is the hash value of the root node.

## Other requirements:

1. Leaf node and Extension node are differentiated by their prefix, not the enum type. The class "Node" is defined in skeleton code and do not change it!

2. General code quality is required. Code smell would be commented for this project, and might affect your grade in the future projects. 


## Hints:

1. Think through all the cases of Get(), Insert(), and Delete() before implementing would save a lot of time. 
2. If you use Rust, implement Clone() for both "Node" and "MerklePatriciaTrie" would be helpful.
3. The grading process includes a lot of test cases, so prepare your own test cases and pass them would help you improve the code. 


