cat raw|awk '{for(i=2;i<=NF;i++){if($i~/https/)print $i}}'|xargs ./get_url|grep -v raw|xargs ./get_vedio|xargs ./download
