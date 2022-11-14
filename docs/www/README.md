---
title: README
category: 636c08cbfcc97f00717832c4
slug: github
hidden: true
---

# Readme.io Documentation

This folder is used to sync documentation with readme.io [here](https://docs.bindplane.observiq.com/docs/about). Each page utilizes a YAML front matter block header that controls the formatting at readme.io. The format for these headers is as follows:

```yaml
---
title: title                        # the title of the page
excerpt: This is a doc page         # the sub-heading excerpt for the page
category: 636c08cbfcc97f00717832c4  # the ID for the category the page is under
parentDoc: 636c0a0cddea2a005d14423a # optional page ID that the doc should be nested under
slug: url                           # the page slug (url)
hidden: false                       # whether the page is published (false) or not (true)
---
```

Category IDs can be found using queries outlined [here](https://docs.readme.com/reference/getcategories). Here is a list of the current categories and their IDs.

| Category | ID |
| --- | --- |
| Getting Started | 636c08cbfcc97f00717832c4 |
| Setup | 636c08d51eb043000f8ce20e |
| Integrations | 636c08e1212e49001e7a3032 |
| Support | 636c08eb959f7600841f16c8 |

Documentation is sync'd using GitHub actions on every tagged release, with new documentation version numbers being set to the BindPlane OP release number.

## Add New Integration

Source, Destination, and Processor pages are nested and require
the `parentDoc` header.

In this example, the category id `636c08e1212e49001e7a3032` points to the Integrations section. The `636c0a0c46142d00a50b384d`
parent doc ID points to the `sources` page.

```
---
title: "Test Source"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "test"
hidden: false
---
```

The following table contains the ID's for all pages which contain
nested documentation

| Page          | ID                         |
| ------------- | -------------------------  |
| Sources       | `636c0a0c46142d00a50b384d` |
| Processors    | `636c0a0cddea2a005d14423a` |
| Destinations  | `636c0a0cd9f114009de9ba78` |
