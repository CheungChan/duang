#[macro_use]
extern crate lazy_static;


//自定义set宏
macro_rules! set {
    ( $( $x:expr ),* ) => {  // Match zero or more comma delimited items
        {
            let mut temp_set = HashSet::new();  // Create a mutable HashSet
            $(
                temp_set.insert($x); // Insert each item matched into the HashSet
            )*
            temp_set // Return the populated HashSet
        }
    };
}

// 自定义hashmap宏
macro_rules! hashmap {
    ($( $key: expr => $val: expr ),*) => {{
         let mut map = std::collections::HashMap::new();
         $( map.insert($key, $val); )*
         map
    }}
}

pub mod scanner;
pub mod ast;
pub mod parser;
pub mod semantic;
pub mod intepretor;