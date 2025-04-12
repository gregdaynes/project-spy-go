/* global sessionStorage */
export const get = (target) => {
  const cachedString = sessionStorage.getItem(target)
  const cachedData = JSON.parse(cachedString)

  return cachedData?.data
}

export const set = (target, data) => {
  const dataString = JSON.stringify({
    data,
  })

  sessionStorage.setItem(target, dataString)

  return get(target)
}

export const remove = (target) => sessionStorage.removeItem(target)

export default {
  get,
  set,
  remove,
}
