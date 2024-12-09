const getImageBase = (repo, repo_name, image_name) => {
  if(repo.repo_type === 'Artifactory' || repo.repo_type === 'Harbor') {
    return `${repo.repo_url.split('://')[1]}/${repo_name}/${image_name}`
  }
  return ""
}

export {
  getImageBase,
}
