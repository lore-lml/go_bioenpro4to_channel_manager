use bioenpro4to_channel_manager::channels::Category;

pub mod utils;
pub mod channel_manager;

fn int_to_category(category: usize) -> anyhow::Result<Category>{
    let res = match category {
        0 => Category::Trucks,
        1 => Category::Scales,
        2 => Category::BioCells,
        _ => return Err(anyhow::Error::msg(""))
    };
    Ok(res)
}
