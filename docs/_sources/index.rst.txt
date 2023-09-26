.. CCNP documentation master file, created by
   sphinx-quickstart on Wed Aug 23 08:18:36 2023.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Welcome to Confidential Cloud-Native Primitives (CCNP)'s documentation!
=======================================================================

VM(Virtual Machine) based confidential computing like Intel TDX provides isolated encryption runtime environment based on
hardware Trusted Execution Environment (TEE) technologies. To land cloud native computing into confidential environment,
there are lots of different PaaS frameworks such as confidential cluster, confidential container, which brings challenges
for enabling and TEE measurement.
This project uses cloud native design pattern to implement confidential computing primitives like event log, measurement,
quote and attestation. It also provides new features design to address new challenges like how to auto scale trustworthy,
how to reduce TCB size, etc.

The project itself contains several parts: the services, the SDK and related dependencies

- Services are designed to hide the complexity of different TEE platforms and provides common interfaces and scalability for cloud-native environment to address the fetching the fetching of quote, measurement and event log.

- SDK is provided to simplify the use of the service interface for development, it covers communication to the service and parses the results from the services. With such SDK, users can perform related actions with one simple API call.

- A ccnp device plugin is provided as the dependency for services such as Quote Server and Measurement Server. It will help with device mount and folder injection within the service.

NOTE: For Intel TDX, it bases on Linux TDX Software Stack at `tdx-tools <https://github.com/intel/tdx-tools>`_, the corresponding white paper is at Whitepaper: `Linux* Stacks for Intel® Trust Domain Extension 1.0 <https://www.intel.com/content/www/us/en/content-details/779108/whitepaper-linux-stacks-for-intel-trust-domain-extension-1-0.html>`_.

.. image:: ccnp_arch.png
   :alt: CCNP architecture
   :align: center


.. toctree::
   :maxdepth: 3
   :caption: Contents:

   CCNP Modules <_rst/modules>



Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
