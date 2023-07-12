use std::env;
use std::path::PathBuf;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::compile_protos("proto/quote_server.proto")?;

    let original_out_dir = PathBuf::from(env::var("OUT_DIR")?);
    let out_dir = "./src";

    tonic_build::configure()
        .out_dir(out_dir)
        .file_descriptor_set_path(original_out_dir.join("quote_server_descriptor.bin"))
        .compile(&["proto/quote_server.proto"], &["proto"])?;

    Ok(())
}
