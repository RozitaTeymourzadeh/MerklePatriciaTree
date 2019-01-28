# cs686_blockchain_P1_Go_skeleton


## Build a Merkle Patricia Tree

Project 1: Build a Merkle Patricia Tree

Ethereum uses a Merkle Patricia Tree (Links to an external site.)Links to an external site. to store the transaction data in a block. By organizing the transaction data in a Merkle Patricia Tree, any block with fraudulent transactions would not match the tree's root hash. Build your own implementation of a Merkle Patricia Trie, following the specifications at the Ethereum wiki. (Links to an external site.)Links to an external site. Be mindful that you will use this code in the subsequent projects. 

Starter Code:

Rust: https://github.com/CharlesGe129/cs686_blockchain_P1_Rust_skeleton.git (Links to an external site.)Links to an external site.

Go: https://github.com/CharlesGe129/cs686_blockchain_P1_Go_skeleton.git (Links to an external site.)Links to an external site.
 

Helpful resources:

Rust: 

https://doc.rust-lang.org/rust-by-example/index.html (Links to an external site.)Links to an external site.

https://doc.rust-lang.org/std/vec/?search= (Links to an external site.)Links to an external site.

Go:

https://tour.golang.org/welcome/1 (Links to an external site.)Links to an external site.

https://golang.org/pkg/strings/#HasPrefix (Links to an external site.)Links to an external site.

 

Project 1 specification

For this project, implement a Merkle Patricia Trie according to this Link (Links to an external site.)Links to an external site. and instructor's lectures. The functions and examples are written in RUST. Skeleton code and some help functions will be provided. You should not change the skeleton code. If you have different ideas that require to change the skeleton code, see the TA *before* submission. Any modification to the skeleton code without the TA's approval will result in point deduction.

You have to choose between Go and Rust. Skeleton code would be provided for both language and please write and pass your own tests before submit the final version of code.

You need to implement five features of the Merkle Patricia Trie:

1. Get()
Description: Get() function takes a key as the argument, traverse down the Merkle Patricia Trie and find the value. If the key doesn't exist, it will return an empty string.
Arguments: key(String)
Return: the value of that key in String type
Rust function definition: get(&mut self, key: &str) -> String

2. Insert()
Description: Insert() function takes a pair of <key, value> as arguments. It will traverse down the Merkle Patricia Trie, find the right place to insert the value, and do the insertion. 
Arguments: key(String), value(String)
Return: None
Rust function definition: get(&mut self, key: &str) -> String

3. Delete()
Description: Delete() function takes a key as the argument, traverse down the Merkle Patricia Trie and find that key. If the key exists, delete the corresponding value and re-balance the trie if necessary; if the key doesn't exist, return "path_not_found".
Arguments: key(String)
Return: None
Rust function definition: delete(&mut self, key: &str)

4. compact_encode()
Description: This function takes an array of HEX value as the input, mark the Node type(such as Branch, Leaf, Extension), make sure the length is even, and convert it into array of ASCII number as the output. You may find the Python version in this Link (Links to an external site.)Links to an external site., but the return type is different!
Arguments: hex_array(array of u8)
Return: array of u8
Rust function definition: compact_encode(hex_array: Vec<u8>) -> Vec<u8>
Example: input=[1, 6, 1], encoded_array=[1, 1, 6, 1], output=[17, 97]

5. compact_decode()
Description: This function reverses the compact_encode() function. 
Arguments: hex_array(array of u8)
Return: array of u8
Rust function definition: compact_decode(encoded_arr: Vec<u8>) -> Vec<u8>
Example: input=[17, 97], output=[1, 6, 1]

Other help functions:

1. fn hash_node(node: &Node) -> String
Description: This function takes a node as the input, hash the node and return the hashed string.

If you use Golang, please follow this link to install the SHA3-256 package: https://github.com/golang/crypto (Links to an external site.)Links to an external site.

If you use Rust, the package dependency is written into Cargo.toml, so directly build and run the project should be fine. 

Classes specification: 
In this project, there are two pre-defined classes. Both classes are defined in skeleton code, feel free to implement any useful functions.

1. enum Node
This class represent a node of type Branch, Leaf, Extension, or Null.

2. struct MerklePatriciaTrie
This class represent a Merkle Patricia Trie. It has two variables: "db" and "root".
Variable "db" is a HashMap. The key of the HashMap is a Node's hash value. The value of the HashMap is the Node. 
Variable "root" is a String, which is the hash value of the root node.

Other requirements:

1. Leaf node and Extension node are differentiated by their prefix, not the enum type. The class "Node" is defined in skeleton code and do not change it!

2. General code quality is required. Code smell would be commented for this project, and might affect your grade in the future projects. 


Hints:

1. Think through all the cases of Get(), Insert(), and Delete() before implementing would save a lot of time. 
2. If you use Rust, implement Clone() for both "Node" and "MerklePatriciaTrie" would be helpful.
3. The grading process includes a lot of test cases, so prepare your own test cases and pass them would help you improve the code. 


