---
title: README
category: 633dd7654359a2003165308a
slug: github
hidden: true
---

# Readme.io Documentation

This folder is used to sync documentation with readme.io [here](https://docs.bindplane.observiq.com/docs/about). Each page utilizes a YAML front matter block header that controls the formatting at readme.io. The format for these headers is as follows:

```yaml
---
title: title                       # the title of the page
excerpt: This is a doc page        # the sub-heading excerpt for the page
category: 633dd7654359a2003165308a # the ID for the category the page is under
slug: url                          # the page slug (url)
hidden: false                      # whether the page is published (false) or not (true)
---
```

Category IDs can be found using queries outlined [here](https://docs.readme.com/reference/getcategories). Here is a list of the current categories and their IDs.

| Category | ID |
| --- | --- |
| Getting Started | 633dd7654359a20031653087 |
| Setup | 633dd7654359a20031653088 |
| Integrations | 633dd7654359a20031653089 |
| Support | 633dd7654359a2003165308a |

Documentation is sync'd using GitHub actions on every tagged release, with new documentation version numbers being set to the BindPlane OP release number.